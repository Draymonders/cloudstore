package handler

import (
	"cloudstore/config"
	mydb "cloudstore/db"
	"cloudstore/meta"
	"cloudstore/mq"
	"cloudstore/store/ceph"
	"cloudstore/store/kodo"
	"cloudstore/util"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// username
// token
// filename
// createtime
// lastedittime
// hash
// filesize

func init() {
	// 目录已经存在
	if _, err := os.Stat(config.DirPath); err == nil {
		return
	}
	// 尝试创建目录
	err := os.MkdirAll(config.DirPath, 0744)
	if err != nil {
		fmt.Println("无法创建临时存储目录，程序将退出")
		os.Exit(1)
	}
}

// getUserName : get username
func getUserName(r *http.Request) string {
	username := r.Form.Get("username")
	// fmt.Println("get username : ", username)
	return username
}

// getToken : get token
func getToken(r *http.Request) string {
	token := r.Form.Get("token")
	// fmt.Println("get token : ", token)
	return token
}

// getFileName : get filename
func getFileName(r *http.Request) string {
	filename := r.Form.Get("filename")
	// fmt.Println("get filename : ", filename)
	return filename
}

// getCreateTime : get create time
func getCreateTime(r *http.Request) string {
	createtime := r.Form.Get("createtime")
	// fmt.Println("get create time:", createtime)
	return createtime
}

// getLastEditTime : get last edit time
func getLastEditTime(r *http.Request) string {
	lastedittime := r.Form.Get("lastedittime")
	// fmt.Println("get last edit time", lastedittime)
	return lastedittime
}

// getHash : get file hash
func getHash(r *http.Request) string {
	hash := r.Form.Get("hash")
	// fmt.Println("get file hash", hash)
	return hash
}

func getFileSize(r *http.Request) int {
	filesize, _ := strconv.Atoi(r.Form.Get("filesize"))
	// fmt.Println("get file size", filesize)
	return filesize
}

// UploadHandler : 上传界面
func UploadHandler(c *gin.Context) {
	c.Redircet("/static/view/index.html")
	return
}

// DoUploadHandler : 需要参数 (username, filename, filesize, filepath, hash)
// DoUploadHandler : 上传文件 (包括秒传接口)
func DoUploadHandler(c *gin.Context) {
	// 1. 从form表单中获得信息
	username := c.Request.FormValue("username")
	file, header, err := r.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "无法获取文件信息",
			"code": util.StatusFormReadError
		})
		return
	}
	defer file.Close()

	filename := header.Filename
	filepath := config.DirPath + filename
	// 2. 创建本地临时文件
	realFile, err := os.Create(filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "无法创建文件",
			"code": util.StatusCreateFileError
		})
		return
	}
	defer realFile.Close()
	// 3. 复制文件，并且计算大小
	filesize, err := io.Copy(realFile, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "无法复制文件到server",
			"code": util.StatusCopyFileError
		})
		return
	}
	// 4. 计算Hash
	realFile.Seek(0, 0)
	hash := util.FileMD5(realFile)

	// 5. 判断之前是否上传过，如果上传过，则只在用户文件表添加记录
	if _, err = meta.IsFileUploadedDB(hash); err == nil {
		suc := mydb.OnUserFileUploadFinished(username,
			filename, hash, int64(filesize))
		if suc {
			resp := util.RespMsg{
				Code: util.StatusOK,
				Msg:  "秒传成功",
			}
		} else {
			resp := util.RespMsg{
				Code: util.StatusFastUploadError,
				Msg:  "秒传失败，请检查数据库数据",
			}
		}
		c.Data(http.StatusOK, "application/json", resp.JSONByte())
		return
	}
	// 6. 记录文件元数据信息
	fileMeta := meta.FileMeta{
		FileName: filename,
		FileSize: filesize,
		FilePath: filepath,
		Hash:     hash,
	}
	
	// 7. 同步或异步将文件转移到oss
	if config.CurrentStoreType == config.StoreCeph {
		log.Println("ceph sync...")
		realFile.Seek(0, 0)
		data, _ := ioutil.ReadAll(realFile)
		cephPath := "/ceph/" + fileMeta.Hash
		err = ceph.PutObject("userfile", cephPath, data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg":  "ceph 上传失败",
				"code": util.StatusCephUploadError
			})
			return
		}
		fileMeta.FilePath = cephPath
		log.Println("ceph end...")
	} else if config.CurrentStoreType == config.StoreKodo {
		fileKey := "userfile/" + fileMeta.FileName
		if !config.AsyncTransferEnable {
			log.Println("kdo sync...")
			// 同步上传到 qiniu kodo
			f := kodo.PutObject(config.KodoBucket, fileMeta.FilePath, fileKey)
			if f == false {
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg":  "kodo 上传失败",
					"code": util.StatusKodoUploadError
				})
				return
			}
			fileMeta.FilePath = kodo.GetObjectURL(fileKey)
		} else {
			fmt.Println("kodo async...")
			data := mq.TransferData{
				FileHash:      fileMeta.Hash,
				CurPath:       fileMeta.FilePath,
				DestPath:      fileKey,
				DestStoreType: config.StoreKodo,
			}
			pubData, _ := json.Marshal(data)
			pubSuc := mq.Publish(
				config.TransExchangeName,
				config.TransKodoRoutingKey,
				pubData,
			)
			if !pubSuc {
				// TODO : 当前转移信息发送失败，稍后重试
			}
		}
		log.Println("kodo end...")
	}
	fmt.Printf("filemeta:{ name:%s path:%s size:%d hash:%s }\n",
		fileMeta.FileName, fileMeta.FilePath, fileMeta.FileSize, fileMeta.Hash)

	// 8. 更新文件表和用户文件表
	flag := meta.CreateFileMetaDB(fileMeta)
	if flag == false {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg" : "更新tb_file表失败"
			"code" : util.StatusStoreToFileError
		})
		return
	}
	flag = mydb.OnUserFileUploadFinished(username, fileMeta.FileName, fileMeta.Hash, fileMeta.FileSize)
	if flag == false {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg" : "更新tb_user_file表失败"
			"code" : util.StatusStoreToUserFileError
		})
		return
	}
	fmt.Printf("UserFile { username:%s filename:%s size:%d hash:%s }\n",
		username, fileMeta.FileName, fileMeta.FileSize, fileMeta.Hash)
	fmt.Println(header.Filename + " uploaded success")
	resp := util.RespMsg{
		Code: util.StatusOK,
		Msg:  "upload success",
	}
	c.Data(http.StatusOK, "Application/json", resp.JSONByte())
}

// GetFileMetaHandler :  查询文件元数据信息
func GetFileMetaHandler(c *gin.Context) {
	filename := c.Request.FormValue("filename")
	fileMeta := meta.GetFileMetaDB(filename)
	fmt.Println("get filename: ", fileName, " fileMeta: ", fileMeta)
	data, _ := json.Marshal(fileMeta)
	c.Data(http.StatusOK, "Application/json", data.JSONByte())
}

// FilesQueryHandler : 文件查询，通过用户名和limit
func FilesQueryHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	token := c.Request.FormValue("token")
	count, _ := strconv.Atoi(c.Request.FormValue("limit")) 
	host := c.Request.Host
	userFiles, err := mydb.QueryUserFileMetas(username, count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg" : "查询 user file表失败"
			"code" : util.StatusQueryUserFilesError
		})
		return
	}
	for i, ufile := range userFiles {
		userFiles[i].DownLoadUrl = DownloadURLHandler(host, username, token, ufile.Hash)
		//  fmt.Printf("username: %s filename: %s\n", ufile.Username, ufile.FileName)
	}
	data, _ := json.Marshal(userFiles)
	c.Data(http.StatusOK, "Application/json", data.JSONByte())
}

// FileDeleteHandler : 文件删除
func FileDeleteHandler(c *gin.Context) {
	filename := c.Request.FormValue("filename")
	fileMeta := meta.GetFileMetaDB(filename)
	// TODO 懒删除，用户删除文件一个星期后删除
	// os.Remove(fileMeta.FilePath)
	meta.RemoveFileMetaDB(filename)
	c.JSON(http.StatusOK, gin.H{
		"msg" : "文件删除成功"
		"code" : util.StatusOK
	})
	return 
}

// FileMetaUpdateHandler : 文件元信息更新
func FileMetaUpdateHandler(c *gin.Context) {
	filename := c.Request.FormValue("filename")
	newfilename :=c.Request.FormValue("newfilename")
	updateSuc := meta.UpdateFileMetaFromfilenameDB(filename, newfilename)
	if !updateSuc {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg" : "文件元数据更新失败"
			"code" : util.StatusFileMetaUpdateError
		})
		return 
	} 
	c.JSON(http.StatusOK, gin.H{
		"msg" : "文件元数据更新成功"
		"code" : util.StatusOK
	})
	return 
}

// TryFastUploadHandler : 秒传接口
func TryFastUploadHandler(c *gin.Context) {
	// 1. 获取参数
	username := c.Request.FormValue("username")
	filename := c.Request.FormValue("filename")
	filesize, _ := strconv.Atoi(c.Request.FormValue("filesize"))
	hash := c.Request.FormValue("hash")
	log.Printf("hash: %s\n", hash)

	// 2. 判断是否在db中
	fileMeta,err := meta.IsFileUploadedDB(hash)
	if fileMeta.Hash == "" {
		fmt.Println(hash, " have not store")
		resp := util.RespMsg{
			Code: util.,
			Msg:  "秒传失败，请访问普通上传接口",
		}
		c.Data(http.StatusBadRequest, "Application/json", resp.JSONByte())
		return
	}
	UserFileCreateSuc = mydb.OnUserFileUploadFinished(username, filename, hash, int64(filesize))
	if UserFileCreateSuc {
		resp := util.RespMsg{
			Code: util.StatusOK,
			Msg:  "秒传成功",
		}
		c.Data(http.StatusOK, "Application/json", resp.JSONByte())
		return
	} else {
		resp := util.RespMsg{
			Code: util.StatusFastUploadError,
			Msg:  "秒传失败，请访问普通上传接口",
		}
		c.Data(http.StatusBadRequest, "Application/json", resp.JSONByte())
		return
	}
}
