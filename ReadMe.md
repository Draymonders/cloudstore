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
- [x] MySQL主从复制
- [x] Redis主从复制 + shell脚本故障转移
- [x] 秒传功能
- [x] 分块上传
- [x] 断点续传
- [ ] Ceph 私有云存储
- [ ] OSS
- [ ] 异步复制 
- [ ] 微服务改造
## 开发环境参数
Go `go version go1.12.5 windows/amd64`
 
操作系统 `Win 10`

IDE `VSCode`

## 文档
- [数据表的建立](./doc/table.sql)
- [MySQL主从同步](./doc/MySQL.md)
- [Redis主从同步](./doc/Redis.md)
- [分块上传原理](./doc/multiPartFileUpload.md)
- [断点续传原理](./doc/BreakpointContinualTransfer.md)