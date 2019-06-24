package main

import (
	. "config"
	mydb "db"
	"encoding/json"
	"fmt"
	"log"
	"mq"
	"store/kodo"
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
	storeSuc := kodo.PutObject(KodoBucket, pubData.CurPath, pubData.DestPath)
	if !storeSuc {
		fmt.Println("上传到Kodo错误，请稍后重试")
		return false
	}
	// 回写数据库
	updateSuc := mydb.UpdateFilePath(pubData.FileHash, pubData.DestPath)
	if !updateSuc {
		fmt.Println("filepath更新失败，请稍后重试")
		return false
	}
	return true
}

func main() {
	if !AsyncTransferEnable {
		log.Println("异步转移文件功能目前被禁用，请检查相关配置")
		return
	}
	log.Println("文件转移服务启动中，开始监听转移任务队列...")
	mq.StartConsume(TransKodoQueueName, "transfer_kodo", ProcessTransfer)
}
