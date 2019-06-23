package main

import (
	"fmt"
	"store/kodo"
)

func main() {
	bucket := "test-kodo"
	localFile := "d:/tmp/irving.png"
	key := "kodo/irving.png"
	f := kodo.PutObject(bucket, localFile, key)
	fmt.Println(f)
	URL := kodo.GetObjectURL(key)
	fmt.Println(URL)
}
