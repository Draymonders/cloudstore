package main

import (
	"fmt"
	"store/ceph"

	"gopkg.in/amz.v1/s3"
)

func main() {
	bucket := ceph.GetCephBucket("userfile")

	err := bucket.PutBucket(s3.PublicRead)
	if err != nil {
		fmt.Println("create bucket error", err.Error())
	}

	// res, err := bucket.List("", "", "", 100)
	// fmt.Printf("object keys: %+v\n", res)

	// err = bucket.Put("/testupload/a.txt", []byte("hello, draymonder"), "octet-stream", s3.PublicRead)
	// fmt.Printf("put object error : %+v\n", err)

	// res, err = bucket.List("", "", "", 100)
	// fmt.Printf("object keys: %+v\n", res)

	// d, _ := bucket.Get("/testupload/a.txt")
	// fmt.Println("data:", string(d))
}
