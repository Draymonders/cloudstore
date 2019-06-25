package main

import (
	"cloudstore/config"
	"cloudstore/route"
)

func main() {
	router := route.Router()
	router.Run(config.UploadServiceHost)
}
