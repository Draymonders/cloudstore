package main

import (
	"fmt"
	"time"

	micro "github.com/micro/go-micro"

	cfg "cloudstore/service/download/config"
	dlProto "cloudstore/service/download/proto"
	"cloudstore/service/download/route"
	dlRpc "cloudstore/service/download/rpc"
)

// startRpcService : 注册服务
func startRPCService() {
	service := micro.NewService(
		micro.Name("go.micro.service.download"), // 在注册中心中的服务名称
		micro.RegisterTTL(time.Second*10),
		micro.RegisterInterval(time.Second*5),
	)
	service.Init()

	dlProto.RegisterDownloadServiceHandler(service.Server(), new(dlRpc.Download))
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

// startApiService : 提供注册服务
func startAPIService() {
	router := route.Router()
	router.Run(cfg.DownloadHost)
}

func main() {
	// api 服务
	go startAPIService()

	// rpc 服务
	startRPCService()
}
