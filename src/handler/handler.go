package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"meta"
	"net/http"
	"os"
	"strconv"
	"time"
	"util"
)

const baseFormat = "2006-01-02 15:04:05"
const dirPath = "/data/tmp/"

// UploadHandler : 上传文件函数
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			fmt.Println("Read file index.html error")
			return
		}
		io.WriteString(w, string(data))
		fmt.Println("Read file index.html")
	} else if r.Method == "POST" {
		// 获取文件
		file, header, err := r.FormFile("file")
		if err != nil {
			fmt.Println("Failed to upload file")
			return
		}
		defer file.Close()
		// 实际存放路径
		filePath := dirPath + header.Filename
		// create fileMeta to store file meta
		fileMeta := meta.FileMeta{
			FileName:   header.Filename,
			FilePath:   filePath,
			CreateTime: time.Now().Format(baseFormat),
		}

		realFile, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Failed to create file: " + header.Filename)
			return
		}
		defer realFile.Close()
		size, err := io.Copy(realFile, file)
		if err != nil {
			fmt.Println("Failed move file to newFile")
			return
		}
		// store size
		fileMeta.FileSize = size
		// store Hash
		realFile.Seek(0, 0)
		fileMeta.Hash = util.FileMD5(realFile)
		fmt.Println(fileMeta)
		// store fileMeta
		flag := meta.CreateFileMetaDB(fileMeta)
		if flag == false {
			fmt.Println(header.Filename + " store meta to db error")
			return
		}
		fmt.Println(header.Filename + " uploaded success")
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

// SucHandler : 上传成功
func SucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload success")
}

// GetFileMetaHandler : 获取文件元数据信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	// 解析Form
	r.ParseForm()

	fileName := r.Form["filename"][0]
	fileMeta := meta.GetFileMetaDB(fileName)
	fmt.Println("get filename: ", fileName, " fileMeta: ", fileMeta)
	data, err := json.Marshal(fileMeta)
	if err != nil {
		fmt.Println("GetFileMetaHandler: parse fileMeta error")
		return
	}
	// tranfer json to client
	w.Write(data)
}

// QueryMultiHandler : 批量获取近期上传的文件
func QueryMultiHandler(w http.ResponseWriter, r *http.Request) {
	// 解析Form
	r.ParseForm()

	count, _ := strconv.Atoi(r.Form.Get("limit"))
	fmt.Println("QueryMultiHandler: count: ", count)
	fileMetas := meta.GetFileMetaListsDB(count)
	data, err := json.Marshal(fileMetas)
	if err != nil {
		fmt.Println("QueryMultiHandler: parse fileMetas error")
		return
	}
	w.Write(data)
}

// FileDownloadHandler : 下载函数
func FileDownloadHandler(w http.ResponseWriter, r *http.Request) {
	// 解析Form
	r.ParseForm()
	filename := r.Form.Get("filename")
	fmt.Println("FileDownloadHandler: filename: ", filename)
	fileMeta := meta.GetFileMetaDB(filename)
	fmt.Println("FileDownloadHandler: fileMeta: ", fileMeta)

	file, err := os.Open(fileMeta.FilePath)
	if err != nil {
		fmt.Println("FileDownloadHandler: can not find the file: ", fileMeta.FilePath)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// close the file
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("FileDownloadHandler: can not read the file: ", fileMeta.FilePath)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// set header
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("content-disposition", "attachment; filename=\""+fileMeta.FileName+"\"")

	// write data to client
	w.Write(data)
}

// FileDeleteHandler : 删除函数
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	// 解析Form
	r.ParseForm()
	filename := r.Form.Get("filename")
	fmt.Println("FileDeleteHandler: filename: ", filename)
	fileMeta := meta.GetFileMetaDB(filename)
	fmt.Println("FileDeleteHandler: fileMeta: ", fileMeta)
	// TODO 懒删除，用户删除文件一个星期后删除
	// os.Remove(fileMeta.FilePath)
	meta.RemoveFileMetaDB(filename)
	fmt.Println("FileDeleteHandler: delete file: ", fileMeta.FilePath, " ok")
	// set ok
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Delete success")
}

// FileMetaUpdateHandler : 文件元信息更改
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// 解析Form
	r.ParseForm()
	filename := r.Form.Get("filename")
	newfilename := r.Form.Get("newfilename")

	fmt.Println("FileMetaUpdateHandler: filename: ", filename, " newfilename: ", newfilename)

	// fileMeta := meta.GetFileMeta(filename)
	// fmt.Println("FileMetaUpdateHandler: old fileMeta: ", fileMeta)
	// remove old fileMeta
	// meta.RemoveFileMeta(filename)
	// store new fileMeta
	// fileMeta.FileName = newfilename
	// meta.UpdateFileMeta(fileMeta)
	// fmt.Println("FileMetaUpdateHandler: new fileMeta: ", fileMeta)
	// set ok
	flag := meta.UpdateFileMetaFromfilenameDB(filename, newfilename)
	if flag == false {
		fmt.Printf("FileMetaUpdateHandler filename %s to %s error", filename, newfilename)
	}
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Update success")
}
