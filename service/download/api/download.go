package api

import (
	"cloudstore/config"
	"cloudstore/db"
	"cloudstore/meta"
	"cloudstore/store/ceph"
	"cloudstore/util"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// DownloadURLHandler : generate the file url
func DownloadURLHandler(c *gin.Context) {
	filehash := c.Request.FormValue("hash")
	fMeta, err := meta.IsFileUploadedDB(filehash)
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code": util.StatusServerError,
				"msg":  "server error",
			})
		return
	}
	if fMeta.FilePath != "" {
		if strings.HasPrefix(fMeta.FilePath, config.DirPath) ||
			strings.HasPrefix(fMeta.FilePath, "/ceph") {
			username := c.Request.FormValue("username")
			token := c.Request.FormValue("token")
			tmpURL := fmt.Sprintf("http://%s/file/download?filehash=%s&username=%s&token=%s",
				c.Request.Host, filehash, username, token)
			c.Data(http.StatusOK, "application/octet-stream", []byte(tmpURL))
		} else if strings.Contains(fMeta.FilePath, "userfile") {
			// oss下载url
			signedURL := fMeta.FilePath
			log.Println(fMeta.FilePath)
			c.Data(http.StatusOK, "application/octet-stream", []byte(signedURL))
		}
	} else {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code": util.StatusServerError,
				"msg":  "server error",
			})
		return
	}
}

// FileDownloadHandler : 下载文件
func FileDownloadHandler(c *gin.Context) {
	hash := c.Request.FormValue("hash")
	username := c.Request.FormValue("username")
	fmt.Println("FileDownloadHandler: hash: ", hash)
	fileMeta, err := meta.IsFileUploadedDB(hash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "通过hash查询file表失败",
			"code": util.StatusQueryFileError,
		})
		return
	}
	userFile, err := db.QueryUserFileMeta(username, hash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "查询user file表失败",
			"code": util.StatusQueryUserFilesError,
		})
		return
	}
	var fileData []byte
	if strings.HasPrefix(fileMeta.FilePath, config.DirPath) {
		fmt.Println("to download file from local...")
		file, err := os.Open(fileMeta.FilePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg":  "文件打开失败",
				"code": util.StatusFileOpenError,
			})
			return
		}
		// close the file
		defer file.Close()
		fileData, err = ioutil.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg":  "文件读取失败",
				"code": util.StatusFileReadError,
			})
			return
		}
	} else if strings.HasPrefix(fileMeta.FilePath, "/ceph") {
		fmt.Println("to download file from ceph...")
		fileData, err = ceph.GetObject("userfile", fileMeta.FilePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg":  "文件读取失败",
				"code": util.StatusFileReadError,
			})
			return
		}
	}
	/*
		else if strings.Contains(fileMeta.FilePath, "userfile") {
			fmt.Println("to download file fom qiniu kodo...")
			resp, err := http.Get("http://" + config.KodoDomain + "/" + fileMeta.FilePath)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg":  "文件读取失败",
					"code": util.StatusFileReadError,
				})
				return
			}
			defer resp.Body.Close()
			fileData, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg":  "文件读取从Resp中读取失败",
					"code": util.StatusFileReadError,
				})
				return
			}
		}
	*/
	// attachment表示文件将会提示下载到本地，而不是直接在浏览器中打开
	c.Header("content-disposition", "attachment; filename=\""+userFile.FileName+"\"")
	// write data to client
	c.Data(http.StatusOK, "application/octect-stream", fileData)
}

// RangeDownloadHandler : download range interface
func RangeDownloadHandler(c *gin.Context) {
	hash := c.Request.FormValue("hash")
	username := c.Request.FormValue("username")
	fileMeta, err := meta.IsFileUploadedDB(hash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "通过hash查询file表失败",
			"code": util.StatusQueryFileError,
		})
		return
	}
	_, err = db.QueryUserFileMeta(username, hash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "查询user file表失败",
			"code": util.StatusQueryUserFilesError,
		})
		return
	}
	f, err := os.Open(fileMeta.FilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "文件打开失败",
			"code": util.StatusFileOpenError,
		})
		return
	}
	defer f.Close()
	// attachment表示文件将会提示下载到本地，而不是直接在浏览器中打开
	c.Header("content-disposition", "attachment; filename=\""+fileMeta.FileName+"\"")
	// write data to client
	http.ServeFile(c.Writer, c.Request, fileMeta.FilePath)
	c.JSON(http.StatusOK, util.NewRespMsg(util.StatusOK, "文件传输完毕", nil).JSONByte())
}
