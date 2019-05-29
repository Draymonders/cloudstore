package main

import (
	"handler"
	"log"
	"net/http"
	"os"
)

func main() {
	// register http 处理函数
	http.HandleFunc("/file/", handler.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
