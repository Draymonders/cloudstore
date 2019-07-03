package main

import (
	"cloudstore/service/apigw/config"
	"cloudstore/service/apigw/route"
)

func main() {
	r := route.Router()
	r.Run(config.UserHost)
}
