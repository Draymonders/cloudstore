package mq

import (
	"config"
	"fmt"

	"github.com/streadway/amqp"
)

var conn *amqp.Connection
var channel *amqp.Channel

var notifyClose chan *amqp.Error

// 初始化，建立可靠Channel
func init() {
	// 要开启异步复制
	if !config.AsyncTransferEnable {
		return
	}
	if initChannel() {
		channel.NotifyClose(notifyClose)
	}
	// 断线自动重连
	go func() {
		for {
			select {
			// 如果接收到notifyClose
			case msg := <-notifyClose:
				conn = nil
				channel = nil
				fmt.Printf("onNotifyChannelClosed: %+v\n", msg)
				initChannel()
			}
		}
	}()
}

// initChannel : 获取 Channel
func initChannel() bool {
	if channel != nil {
		return true
	}

	conn, err := amqp.Dial(config.RabbitURL)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	channel, err = conn.Channel()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

// Publish : 往exchange 发送 msg 路由绑定 routingKey
func Publish(exchange, routingKey string, msg []byte) bool {
	if !initChannel() {
		return false
	}
	if nil == channel.Publish(exchange, routingKey, false, // 如果没有对应的queue, 就会丢弃这条消息
		false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg}) {
		return true
	}
	return false
}
