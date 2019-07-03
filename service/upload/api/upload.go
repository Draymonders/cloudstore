package api

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
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

// DoUploadHandler : 需要参数 (username, filename, filesize, filepath, hash)
// DoUploadHandler : 上传文件 (包括秒传接口)
func DoUploadHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	// 1. 从form表单中获得信息
	username := c.Request.FormValue("username")
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "无法获取文件信息",
			"code": util.StatusFormReadError,
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
			"code": util.StatusCreateFileError,
		})
		return
	}
	defer realFile.Close()
	// 3. 复制文件，并且计算大小
	filesize, err := io.Copy(realFile, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "无法复制文件到server",
			"code": util.StatusCopyFileError,
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
			c.Data(http.StatusOK, "application/json", resp.JSONByte())
		} else {
			resp := util.RespMsg{
				Code: util.StatusFastUploadError,
				Msg: "秒传失败，请检查	数据库数据",
			}
			c.Data(http.StatusInternalServerError, "application/json", resp.JSONByte())
		}
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
				"code": util.StatusCephUploadError,
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
					"code": util.StatusKodoUploadError,
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
			"msg":  "更新tb_file表失败",
			"code": util.StatusStoreToFileError,
		})
		return
	}
	flag = mydb.OnUserFileUploadFinished(username, fileMeta.FileName, fileMeta.Hash, fileMeta.FileSize)
	if flag == false {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "更新tb_user_file表失败",
			"code": util.StatusStoreToUserFileError,
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

// TryFastUploadHandler : 秒传接口
func TryFastUploadHandler(c *gin.Context) {
	// 1. 获取参数
	username := c.Request.FormValue("username")
	filename := c.Request.FormValue("filename")
	filesize, _ := strconv.Atoi(c.Request.FormValue("filesize"))
	hash := c.Request.FormValue("hash")
	log.Printf("hash: %s\n", hash)

	// 2. 判断是否在db中
	fileMeta, _ := meta.IsFileUploadedDB(hash)
	if fileMeta.Hash == "" {
		fmt.Println(hash, " have not store")
		resp := util.RespMsg{
			Code: util.StatusQueryFileError,
			Msg:  "秒传失败，请访问普通上传接口",
		}
		c.Data(http.StatusBadRequest, "Application/json", resp.JSONByte())
		return
	}
	UserFileCreateSuc := mydb.OnUserFileUploadFinished(username, filename, hash, int64(filesize))
	if UserFileCreateSuc {
		resp := util.RespMsg{
			Code: util.StatusOK,
			Msg:  "秒传成功",
		}
		c.Data(http.StatusOK, "Application/json", resp.JSONByte())
		return
	}
	resp := util.RespMsg{
		Code: util.StatusFastUploadError,
		Msg:  "秒传失败，请访问普通上传接口",
	}
	c.Data(http.StatusBadRequest, "Application/json", resp.JSONByte())
	return
}
