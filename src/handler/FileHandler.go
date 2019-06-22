package handler

import (
	mydb "db"
	"encoding/json"
	"fmt"
	"io"
	"meta"
	"net/http"
	"os"
	"strconv"
	"util"
)

// getUserName : get username
func getUserName(r *http.Request) string {
	username := r.Form.Get("username")
	fmt.Println("get username : ", username)
	return username
}

// getToken : get token
func getToken(r *http.Request) string {
	token := r.Form.Get("token")
	fmt.Println("get token : ", token)
	return token
}

// getFileName : get filename
func getFileName(r *http.Request) string {
	filename := r.Form.Get("filename")
	fmt.Println("get filename : ", filename)
	return filename
}

// getCreateTime : get create time
func getCreateTime(r *http.Request) string {
	createtime := r.Form.Get("createtime")
	fmt.Println("get create time:", createtime)
	return createtime
}

// getLastEditTime : get last edit time
func getLastEditTime(r *http.Request) string {
	lastedittime := r.Form.Get("lastedittime")
	fmt.Println("get last edit time", lastedittime)
	return lastedittime
}

// getHash : get file hash
func getHash(r *http.Request) string {
	hash := r.Form.Get("hash")
	fmt.Println("get file hash", hash)
	return hash
}

func getFileSize(r *http.Request) int {
	filesize, _ := strconv.Atoi(r.Form.Get("filesize"))
	fmt.Println("get file size", filesize)
	return filesize
}

// StatusInternalServerError : set w head
func StatusInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}

// UploadHandler : need param (username, filename, filesize, filepath, hash)
// UploadHandler : upload file (include fast upload)
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// read file
		http.Redirect(w, r, "/static/view/index.html", http.StatusFound)
		fmt.Println("Read file index.html")
		return
	} else if r.Method == "POST" {
		// parse data from form
		r.ParseForm()
		username := getUserName(r)
		// get the file from form
		file, header, err := r.FormFile("file")
		if err != nil {
			fmt.Println("Failed to read file from form, err:", err.Error())
			StatusInternalServerError(w)
			return
		}
		defer file.Close()
		// real Path
		filename := header.Filename
		filepath := dirPath + filename
		realFile, err := os.Create(filepath)
		if err != nil {
			fmt.Println("Failed to create file: " + header.Filename)
			StatusInternalServerError(w)
			return
		}
		defer realFile.Close()
		filesize, err := io.Copy(realFile, file)
		if err != nil {
			fmt.Println("Failed move file to newFile, err", err.Error())
			StatusInternalServerError(w)
			return
		}

		// TODO split to Hash service

		// store Hash
		realFile.Seek(0, 0)
		hash := util.FileMD5(realFile)

		if _, err = meta.IsFileUploadedDB(hash); err == nil {
			// file has store to local
			// then only store data to tb_user_file
			suc := mydb.OnUserFileUploadFinished(username,
				filename, hash, int64(filesize))
			if suc {
				resp := util.RespMsg{
					Code: 0,
					Msg:  "秒传成功",
				}
				fmt.Println(filename, "秒传成功")
				w.Write(resp.JSONByte())
				return
			}
			resp := util.RespMsg{
				Code: -1,
				Msg:  "秒传失败，请检查数据库数据",
			}
			StatusInternalServerError(w)
			w.Write(resp.JSONByte())
			return
		}

		fileMeta := meta.FileMeta{
			FileName: filename,
			FileSize: filesize,
			FilePath: filepath,
			Hash:     hash,
		}
		fmt.Printf("filemeta:{ name:%s path:%s size:%d hash:%s }\n",
			fileMeta.FileName, fileMeta.FilePath, fileMeta.FileSize, fileMeta.Hash)
		// store fileMeta
		flag := meta.CreateFileMetaDB(fileMeta)
		if flag == false {
			fmt.Println(header.Filename + " store meta to file db error")
			StatusInternalServerError(w)
			return
		}
		flag = mydb.OnUserFileUploadFinished(username, fileMeta.FileName, fileMeta.Hash, fileMeta.FileSize)
		if flag == false {
			fmt.Println(header.Filename + " store meta to user file db error")
			StatusInternalServerError(w)
			return
		}
		fmt.Printf("UserFile { username:%s filename:%s size:%d hash:%s }\n",
			username, fileMeta.FileName, fileMeta.FileSize, fileMeta.Hash)
		fmt.Println(header.Filename + " uploaded success")
		resp := util.RespMsg{
			Code: 0,
			Msg:  "upload success",
		}
		w.Write(resp.JSONByte())
	}
}

// GetFileMetaHandler : query more about a file
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	// 解析Form
	r.ParseForm()

	fileName := getFileName(r)
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
	username := getUserName(r)
	count, _ := strconv.Atoi(r.Form.Get("limit"))
	fmt.Println("QueryMultiHandler: limit: ", count)
	// fileMetas := meta.(count)
	userFiles, err := mydb.QueryUserFileMetas(username, count)
	if err != nil {
		fmt.Println("QueryMultiHandler err:", err.Error())
		return
	}
	for i, ufile := range userFiles {
		userFiles[i].DownLoadUrl = DownloadURLHandler(r, ufile.Hash)
		fmt.Printf("username:%s filename:%s size:%d hash:%s url:%s\n", ufile.Username, ufile.FileName, ufile.FileSize,
			ufile.Hash, ufile.DownLoadUrl)
	}

	data, err := json.Marshal(userFiles)
	if err != nil {
		fmt.Println("QueryMultiHandler: parse userFiles error")
		return
	}
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
	hash := r.Form.Get("hash")

	fmt.Printf("hash: %s\n", hash)
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

/*
if _, err = meta.IsFileUploadedDB(hash); err == nil {
			// file has store to local
			// then only store data to tb_user_file
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
				Code: -1,
				Msg:  "秒传失败，请检查数据库数据",
			}
			StatusInternalServerError(w)
			w.Write(resp.JSONByte())
			return
		}
*/
