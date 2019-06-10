package handler

import (
	mydb "db"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"meta"
	"net/http"
	"os"
	"strconv"
	"util"
)

const baseFormat = "2006-01-02 15:04:05"
const dirPath = "/data/tmp/"

// UploadHandler : 上传文件函数
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// data, err := ioutil.ReadFile("./static/view/index.html")
		// if err != nil {
		// 	fmt.Println("Read file index.html error")
		// 	return
		// }
		// io.WriteString(w, string(data))
		http.Redirect(w, r, "/static/view/index.html", http.StatusFound)
		fmt.Println("Read file index.html")
		return
	} else if r.Method == "POST" {
		// parse data from form
		r.ParseForm()

		// get username
		username := r.Form.Get("username")

		// get the file from buffer
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
			FileName: header.Filename,
			FilePath: filePath,
		}
		// new a real file to store file
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
		// TODO split to Hash service

		// store Hash
		realFile.Seek(0, 0)
		fileMeta.Hash = util.FileMD5(realFile)
		fmt.Println(fileMeta)

		// store fileMeta
		flag := meta.CreateFileMetaDB(fileMeta)
		if flag == false {
			fmt.Println(header.Filename + " store meta to file db error")
			return
		}
		flag = mydb.OnUserFileUploadFinished(username, fileMeta.FileName, fileMeta.Hash, fileMeta.FileSize)
		if flag == false {
			fmt.Println(header.Filename + " store meta to user file db error")
			return
		}
		fmt.Println(header.Filename + " uploaded success")
		http.Redirect(w, r, "/static/view/home.html", http.StatusFound)
	}
}

// SucHandler : upload success
func SucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload success")
}

// GetFileMetaHandler : query more about a file
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

// QueryMultiHandler : query user files for(username, limit)
func QueryMultiHandler(w http.ResponseWriter, r *http.Request) {
	// 解析Form
	r.ParseForm()
	username := r.Form.Get("username")
	count, _ := strconv.Atoi(r.Form.Get("limit"))
	fmt.Println("QueryMultiHandler: count: ", count)
	// fileMetas := meta.(count)
	userFiles, err := mydb.QueryUserFileMetas(username, count)
	if err != nil {
		fmt.Println("QueryMultiHandler err:", err.Error())
		return
	}
	data, err := json.Marshal(userFiles)
	if err != nil {
		fmt.Println("QueryMultiHandler: parse userFiles error")
		return
	}
	w.Write(data)
}

// FileDownloadHandler : download file
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

// FileDeleteHandler : remove file meta
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

// FileMetaUpdateHandler : update file meta
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

// TryFastUploadHandler : fast upload handler
func TryFastUploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// 1. parse req
	username := r.Form.Get("username")
	filename := r.Form.Get("filename")
	filesize, _ := strconv.Atoi(r.Form.Get("filesize"))
	hash := r.Form.Get("filehash")

	// 2. check if is same exists
	fileMeta, err := meta.IsFileUploadedDB(hash)
	if err != nil {
		fmt.Println("TryFastUploadHandler : err: ", err.Error())
		return
	}
	if fileMeta.Hash == "" {
		fmt.Println(hash, " have not store")
		resp := util.RespMsg{
			Code: -1,
			Msg:  "秒传失败，请访问普通上传接口",
		}
		w.Write(resp.JSONByte())
		return
	}
	suc := mydb.OnUserFileUploadFinished(username,
		filename, hash, int64(filesize))
	if suc {
		resp := util.RespMsg{
			Code: 0,
			Msg:  "秒传成功",
		}
		w.Write(resp.JSONByte())
		return
	}
	resp := util.RespMsg{
		Code: -2,
		Msg:  "秒传失败，请稍后重试",
	}
	w.Write(resp.JSONByte())
	return
}
