package handler

import (
	. "config"
	"db"
	"fmt"
	"io/ioutil"
	"meta"
	"net/http"
	"os"
	"store/ceph"
	"strings"
	"util"
)

// DownloadURLHandler : generate the file url
func DownloadURLHandler(r *http.Request, filehash string) string {
	// filehash := getHash(r)
	username := getUserName(r)
	token := getToken(r)
	tmpURL := fmt.Sprintf(
		"http://%s/file/download?hash=%s&username=%s&token=%s", r.Host, filehash, username, token)
	// w.Write([]byte(tmpURL))
	return tmpURL
}

// FileDownloadHandler : download file
func FileDownloadHandler(w http.ResponseWriter, r *http.Request) {
	// 解析Form
	r.ParseForm()
	hash := getHash(r)
	username := getUserName(r)
	fmt.Println("FileDownloadHandler: hash: ", hash)
	fileMeta, err := meta.IsFileUploadedDB(hash)
	if err != nil {
		fmt.Println("FileDownloadHandler: query file hash failed, err :", err.Error())
		StatusInternalServerError(w)
		w.Write(util.NewRespMsg(-1, "query file through hash failed", nil).JSONByte())
		return
	}
	userFile, err := db.QueryUserFileMeta(username, hash)
	if err != nil {
		fmt.Println("FileDownloadHandler: query userfile failed, err :", err.Error())
		StatusInternalServerError(w)
		w.Write(util.NewRespMsg(-1, "get UserFile through hash and username failed", nil).JSONByte())
		return
	}
	var fileData []byte
	if strings.HasPrefix(fileMeta.FilePath, DirPath) {
		fmt.Println("to download file from local...")
		file, err := os.Open(fileMeta.FilePath)
		if err != nil {
			fmt.Println("FileDownloadHandler: can not find the file: ", fileMeta.FilePath)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// close the file
		defer file.Close()
		fileData, err = ioutil.ReadAll(file)
		if err != nil {
			fmt.Println("FileDownloadHandler: can not read data from file: ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if strings.HasPrefix(fileMeta.FilePath, "/ceph") {
		fmt.Println("to download file from ceph...")
		fileData, err = ceph.GetObject("userfile", fileMeta.FilePath)
		if err != nil {
			fmt.Println("FileDownloadHandler: can not read data from file: ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// set header
	w.Header().Set("Content-Type", "application/octect-stream")
	// attachment表示文件将会提示下载到本地，而不是直接在浏览器中打开
	w.Header().Set("content-disposition", "attachment; filename=\""+userFile.FileName+"\"")

	// write data to client
	w.Write(fileData)
}

// RangeDownloadHandler : download range interface
func RangeDownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	hash := getHash(r)
	username := getUserName(r)
	fileMeta, err := meta.IsFileUploadedDB(hash)
	if err != nil {
		fmt.Println("RangeDownloadHandler: query file hash failed, err :", err.Error())
		StatusInternalServerError(w)
		w.Write(util.NewRespMsg(-1, "query file through hash failed", nil).JSONByte())
		return
	}
	userfile, err := db.QueryUserFileMeta(username, hash)
	if err != nil {
		fmt.Println("RangeDownloadHandler: query userfile failed, err :", err.Error())
		StatusInternalServerError(w)
		w.Write(util.NewRespMsg(-1, "get UserFile through hash and username failed", nil).JSONByte())
		return
	}
	f, err := os.Open(fileMeta.FilePath)
	if err != nil {
		fmt.Println("RangeDownloadHandler: can not find the file: ", fileMeta.FilePath)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	w.Header().Set("Content-Type", "application/octect-stream")
	// attachment表示文件将会提示下载到本地，而不是直接在浏览器中打开
	w.Header().Set("content-disposition", "attachment; filename=\""+userfile.FileName+"\"")
	http.ServeFile(w, r, fileMeta.FilePath)
}
