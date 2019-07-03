package main

import (
	"cloudstore/service/apigw/route"
)

func main() {
	r := route.Router()
	r.Run(":8080")
}
