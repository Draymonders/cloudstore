package main

import (
	"cloudstore/config"
	"cloudstore/mq"
	"cloudstore/service/upload/route"
	"fmt"
	"log"
	"time"

	upProto "cloudstore/service/upload/proto"
	upRpc "cloudstore/service/upload/rpc"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
)

// startUploadService : 开始上传服务
func startRPCService() {
	service := micro.NewService(
		micro.Name("go.micro.service.upload"),
		micro.RegisterTTL(time.Second*10),     // TTL指定从上一次心跳间隔起，超过这个时间服务会被服务发现移除
		micro.RegisterInterval(time.Second*5), // 让服务在指定时间内重新注册，保持TTL获取的注册时间有效
		// micro.Flags(common.CustomFlags...),
	)
	service.Init(
		micro.Action(func(c *cli.Context) {
			// 检查是否指定mqhost
			mqhost := c.String("mqhost")
			if len(mqhost) > 0 {
				log.Println("custom mq address: " + mqhost)
				mq.UpdateRabbitHost(mqhost)
			}
		}),
	)
	upProto.RegisterUploadServiceHandler(service.Server(), new(upRpc.Upload))
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

func startAPIService() {
	router := route.Router()
	router.Run(config.UploadHost)
}

func main() {
	// 启动内部服务
	go startAPIService()
	// 等待RPC连接
	startRPCService()
}
