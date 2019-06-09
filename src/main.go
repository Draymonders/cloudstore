package main

import (
	"fmt"
	"handler"
	"net/http"
)

const listenPort = ":80"

func main() {
	// 静态资源处理
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	// 动态路由
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.SucHandler)
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/file/query", handler.QueryMultiHandler)
	http.HandleFunc("/file/down", handler.FileDownloadHandler)
	http.HandleFunc("/file/del", handler.FileDeleteHandler)
	http.HandleFunc("/file/update", handler.FileMetaUpdateHandler)

	http.HandleFunc("/user/signup", handler.SignupHandler)
	http.HandleFunc("/user/signin", handler.SignInHandler)
	http.HandleFunc("/user/info", handler.HTTPInterceptor(handler.UserInfoHandler))
	fmt.Println("server start , listen", listenPort)
	err := http.ListenAndServe(listenPort, nil)
	if err != nil {
		fmt.Println("server start failed err:" + err.Error() + "\n")
		return
	}
}
