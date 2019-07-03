package main

import (
	"cloudstore/config"
	"cloudstore/service/apigw/route"
)

func main() {
	r := route.Router()
	r.Run(config.UserHost)
}
