package api

import (
	rPool "cloudstore/cache/redis"
	"cloudstore/config"
	mydb "cloudstore/db"
	"cloudstore/util"
	"fmt"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
)

const chunkSize = 5 * 1024 * 1024

// MultipartUploadInfo : multi part upload info
type MultipartUploadInfo struct {
	UploadID   string
	FileSize   int
	FileHash   string
	ChunkSize  int
	ChunkCount int
}

// InitialMultipartUploadHandler : 初始化分块信息
func InitialMultipartUploadHandler(c *gin.Context) {
	// 1. 获取参数
	username := c.Request.FormValue("username")
	filehash := c.Request.FormValue("hash")
	filesize, _ := strconv.Atoi(c.Request.FormValue("filesize"))

	// 2. 获取Redis链接
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// 3. 创建分块信息
	upInfo := MultipartUploadInfo{
		UploadID:   username + fmt.Sprintf("%x", time.Now().Unix()),
		FileSize:   filesize,
		FileHash:   filehash,
		ChunkSize:  chunkSize,
		ChunkCount: int(math.Ceil(float64(filesize) / float64(chunkSize))),
	}

	// 4. 存储到Redis中
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "chunkcount", upInfo.ChunkCount)
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "filehash", upInfo.FileHash)
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "filesize", upInfo.FileSize)

	// 5. 返回文件
	resp := util.NewRespMsg(util.StatusOK, "分块信息初始化成功", nil)
	c.Data(http.StatusOK, "application/json", resp.JSONByte())
}

// UploadPartHandler : 分块信息上传记录
func UploadPartHandler(c *gin.Context) {
	// 1. 获取上传的分块标号
	uploadID := c.Request.FormValue("uploadid")
	chunkIndex := c.Request.FormValue("index")

	// 2. 获取Redis链接
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// 3. 保存文件分块到server本地
	fpath := config.DirPath + uploadID + "/" + chunkIndex
	os.MkdirAll(path.Dir(fpath), 0744)
	fd, err := os.Create(fpath)
	if err != nil {
		// StatusCreateFileError : 7 创建文件失败
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "创建文件失败",
			"code": util.StatusCreateFileError,
		})
		return
	}
	defer fd.Close()

	// 4. 将Form的data传输到server中
	buf := make([]byte, 1024*1024)
	for {
		n, err := c.Request.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil {
			break
		}
	}

	// 5. 将已经读取的分块信息保存到Redis的Hash Object中  key为 MP_uploadID
	rConn.Do("HSET", "MP_"+uploadID, "chkidx_"+chunkIndex, 1)
	resp := util.NewRespMsg(util.StatusOK, "分块信息"+chunkIndex+"上传成功", nil)
	c.Data(http.StatusOK, "application/json", resp.JSONByte())
}

// CompleteUploadHandler : 分块完成查询成功与否
func CompleteUploadHandler(c *gin.Context) {
	// 1. 获取参数
	username := c.Request.FormValue("username")
	filename := c.Request.FormValue("filename")
	filehash := c.Request.FormValue("hash")
	filesize, _ := strconv.Atoi(c.Request.FormValue("filesize"))
	uploadID := c.Request.FormValue("uploadid")

	// 2. 获取Redis链接
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// 3. 使用HGETALL MP_uploadid, 获取hash object中key为MP_uploadid的所有value
	data, err := redis.Values(rConn.Do("HGETALL", "MP_"+uploadID))
	if err != nil {
		// StatusRedisGetError : 20 redis获取数据失败
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": util.StatusRedisGetError,
			"msg":  "redis获取数据失败",
		})
		return
	}

	// 4. 读取上面数据，判断所有分块是否已经上传完毕
	totalCount := 0
	chunkCount := 0
	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))
		if k == "chunkcount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			chunkCount++
		}
	}

	// 5. 判断分块是否上传成功
	if totalCount != chunkCount {
		// StatusPartNotFull : 21 分块没有上传完毕
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "分块没有上传完毕",
			"code": util.StatusPartNotFull,
		})
		return
	}

	// 6. 合并分块为一个单独的文件，本质上使用的是如下命令
	// cat `ls | sort -n` > /tmp/filename

	partFileStorePath := config.DirPath + "/uploadid" // 分块所在的目录
	fileStorePath := config.DirPath + filename        // 最后文件保存的路径
	if _, err := mergeAllPartFile(partFileStorePath, fileStorePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "分块归并失败",
			"code": util.StatusMergeError,
		})
		return
	}

	// 7. 将文件信息写回到mysql
	mydb.OnFileUploadFinished(filename, int64(filesize), fileStorePath, filehash)
	mydb.OnUserFileUploadFinished(username, filename, filehash, int64(filesize))

	resp := util.NewRespMsg(util.StatusOK, "分块处理完毕", nil)
	c.Data(http.StatusOK, "application/json", resp.JSONByte())
}

// mergeAllPartFile: filepath： 分块存储的路径 filestore： 文件最终地址
func mergeAllPartFile(partFileStorePath, fileStorePath string) (bool, error) {
	var cmd *exec.Cmd
	cmd = exec.Command(config.MergeAllShell, partFileStorePath, fileStorePath)

	if _, err := cmd.Output(); err != nil {
		fmt.Println(err)
		return false, err
	}
	fmt.Println(fileStorePath, " has been merge complete")
	return true, nil
}
