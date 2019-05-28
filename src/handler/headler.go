package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// 上传处理
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// 返回上传html页面
		data, e := ioutil.ReadFile("./static/view/index.html")
		if e != nil {
			io.WriteString(w, "internel server error")
			return
		}
		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		// 接收文件流及存储到本地目录
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to get data, err: %s\n", err.Error())
			return
		}
		defer file.Close()

		// 创建本地文件 接受文件流
		localFile, err := os.Create("/tmp/"+head.Filename)
		if err != nil {
			fmt.Printf("Failed to create file, err: %s\n", err.Error())
			return
		}
		defer localFile.Close()

		_, err = io.Copy(localFile, file)
		if err != nil {
			fmt.Printf("Failed to save data into file, err:%s", err.Error())
			return
		}
		http.Redirect(w,r,"/file/upload/suc", http.StatusFound)
	}
	return
}


func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload succeed!")
}