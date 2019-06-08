package main

import (
	"fmt"
	"handler"
	"net/http"
)

func main() {
	const listenPort = ":80"

	// 动态路由
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.SucHandler)
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/file/query", handler.QueryMultiHandler)
	http.HandleFunc("/file/down", handler.FileDownloadHandler)
	http.HandleFunc("/file/del", handler.FileDeleteHandler)
	http.HandleFunc("/file/update", handler.FileMetaUpdateHandler)
	// http.HandleFunc("/file/")
	err := http.ListenAndServe(listenPort, nil)
	if err != nil {
		fmt.Println("server start failed err:" + err.Error() + "\n")
		return
	}

}
