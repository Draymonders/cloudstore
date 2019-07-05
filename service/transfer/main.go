package main

import (
	"cloudstore/config"
	cfg "cloudstore/config"
	mydb "cloudstore/db"
	"cloudstore/mq"
	"cloudstore/store/kodo"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/micro/go-micro"
)

// ProcessTransfer : 处理接收到的队列里面的数据
func ProcessTransfer(msg []byte) bool {
	log.Println(string(msg))

	pubData := mq.TransferData{}
	err := json.Unmarshal(msg, &pubData)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	storeSuc := kodo.PutObject(cfg.KodoBucket, pubData.CurPath, pubData.DestPath)
	if !storeSuc {
		fmt.Println("上传到Kodo错误，请稍后重试")
		return false
	}
	destPath := kodo.GetObjectURL(pubData.DestPath)
	// 回写数据库
	updateSuc := mydb.UpdateFilePath(pubData.FileHash, destPath)
	if !updateSuc {
		fmt.Println("filepath更新失败，请稍后重试")
		return false
	}
	return true
}

// startRPCService : 注册转移服务
func startRPCService() {
	service := micro.NewService(
		micro.Name("go.micro.service.transfer"), // 服务名称
		micro.RegisterTTL(time.Second*10),       // TTL指定从上一次心跳间隔起，超过这个时间服务会被服务发现移除
		micro.RegisterInterval(time.Second*5),   // 让服务在指定时间内重新注册，保持TTL获取的注册时间有效
	)
	service.Init()

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

// startTranserService ： 注册转移服务
func startTranserService() {
	if !config.AsyncTransferEnable {
		log.Println("异步转移文件功能目前被禁用，请检查相关配置")
		return
	}
	log.Println("文件转移服务启动中，开始监听转移任务队列...")
	mq.StartConsume(
		config.TransKodoQueueName,
		"transfer_kodo",
		ProcessTransfer)
}

func main() {
	// 文件转移服务
	go startTranserService()

	// rpc 服务
	startRPCService()
}
