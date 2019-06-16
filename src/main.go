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

	// file upload (include fast upload)
	http.HandleFunc("/file/upload", handler.HTTPInterceptor(handler.UploadHandler))
	http.HandleFunc("/file/meta", handler.HTTPInterceptor(handler.GetFileMetaHandler))
	http.HandleFunc("/file/query", handler.HTTPInterceptor(handler.QueryMultiHandler))
	http.HandleFunc("/file/download", handler.HTTPInterceptor(handler.FileDownloadHandler))
	http.HandleFunc("/file/del", handler.HTTPInterceptor(handler.FileDeleteHandler))
	http.HandleFunc("/file/update", handler.HTTPInterceptor(handler.FileMetaUpdateHandler))

	// file down range
	http.HandleFunc("/file/downloadurl", handler.HTTPInterceptor(DownloadURLHandler))
	http.HandleFunc("/file/download/range", handler.HTTPInterceptor(RangeDownloadHandler))

	// fast upload
	// http.HandleFunc("/file/fastupload", handler.TryFastUploadHandler)

	// mutli part file upload
	http.HandleFunc("/file/mpupload/init", handler.HTTPInterceptor(
		handler.InitialMultipartUploadHandler))
	http.HandleFunc("/file/mpupload/uppart", handler.HTTPInterceptor(
		handler.UploadPartHandler))
	http.HandleFunc("/file/mpupload/complete", handler.HTTPInterceptor(
		handler.CompleteUploadHandler))

	// user vertify
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
