package handler

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

//  handle request
func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m == http.MethodPut {
		put(w, r)
		return
	} else if m == http.MethodGet {
		get(w, r)
		return
	}
	// get error , return @code 405 Method Not Allowed
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func put(w http.ResponseWriter, r *http.Request) {
	// EscapedPath
	filePath := "/file/" + strings.Split(r.URL.EscapedPath(), "/")[2]
	log.Println("put path: ", filePath)
	// now f is a io.Writer
	f, e := os.Create(os.Getenv("STORE_ROOT") + filePath)
	if e != nil {
		log.Println(e)
		// get error , return @code 500 Internal Server Error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	io.Copy(f, r.Body)
}

func get(w http.ResponseWriter, r *http.Request) {
	filePath := "/file/" + strings.Split(r.URL.EscapedPath(), "/")[2]
	log.Println("get path: ", filePath)
	// now f is a io.Reader
	f, e := os.Open(os.Getenv("STORE_ROOT") + filePath)
	if e != nil {
		log.Println(e)
		// get error, return @code 404 Status Not Found
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	io.Copy(w, f)
}
