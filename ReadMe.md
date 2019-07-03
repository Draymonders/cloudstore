# 基于golang实现的一种简易分布式云存储服务

linux默认存储路径`/data/tmp/`

启动服务
```shell 
# windows 启动 consul
consul agent -dev
# 启动 gateway 服务
go run service/apigw/main.go --registry=consul
# 启动 account 服务
go run service/account/main.go --registry=consul
# 启动 upload 服务
go run service/upload/main.go --registry=consul
# 启动 download 服务
go run service/download/main.go --registry=consul
```
## 功能
- [x] 单机文件存储
- [x] MySQL 主从复制
- [x] Redis 主从复制 + shell 脚本故障转移
- [x] 秒传功能
- [x] 分块上传
- [x] ~~断点续传~~
- [x] Ceph 私有云存储
- [x] Kodo 公有云存储 (七牛云对象存储)
- [x] Rabbitmq 异步复制 
- [x] 微服务改造
- [] 运维自动化
## 开发环境参数


操作系统 `Win 10`&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;IDE `VSCode`&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Go `go version go1.12.5 windows/amd64`
 
## 文档
- [Go 入门](https://tour.go-zh.org/welcome/1)
- [MySQL 使用手册](https://chhy2009.github.io/document/mysql-reference-manual.pdf)
- [数据表的建立](./doc/table.sql)
- [MySQL 主从同步](./doc/MySQL.md)
- [Redis 命令手册](http://redisdoc.com/)
- [Redis 主从同步](./doc/Redis.md)
- [分块上传原理](./doc/multiPartFileUpload.md)
- [断点续传原理](./doc/BreakpointContinualTransfer.md)
- [Ceph 中文社区](http://ceph.org.cn/) [Ceph 中文文档](http://docs.ceph.org.cn/)
- [Ceph 私有云存储实践](./doc/ceph.md)
- [七牛云 kodo 使用体验](./doc/kodo.md)
- [RabbitMQ 使用体验](./doc/rabbitmq.md) 
- [RabbitMQ 英文官方](http://www.rabbitmq.com/getstarted.html) [RabbitMQ 一个中文版文档](http://rabbitmq.mr-ping.com/)
- [gRPC 官方文档中文版](http://doc.oschina.net/grpc?t=56831)
- [go-micro微服务框架 github源码](https://github.com/micro/go-micro)
- [gin web框架 github源码](https://github.com/gin-gonic/gin)
- [k8s 中文社区](https://www.kubernetes.org.cn/docs)