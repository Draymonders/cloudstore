# 简易分布式对象存储 -- GO实现

默认启动端口`80`

linux默认存储路径`/data/tmp/`

windows默认存储路径`D:\tmp\`

启动方式

```
go run main.go
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
- [ ] 微服务改造
## 开发环境参数


操作系统 `Win 10`&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;IDE `VSCode`&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Go `go version go1.12.5 windows/amd64`
 


## 文档
- [数据表的建立](./doc/table.sql)
- [MySQL 主从同步](./doc/MySQL.md)
- [Redis 主从同步](./doc/Redis.md)
- [分块上传原理](./doc/multiPartFileUpload.md)
- [断点续传原理](./doc/BreakpointContinualTransfer.md)
- [Ceph 私有云存储实践](./doc/ceph.md)
- [七牛云 kodo 使用体验](./doc/kodo.md)