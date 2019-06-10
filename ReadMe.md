# 简易分布式对象存储 -- GO实现

目前已实现简易单机版存储 + MySQL (未用Redis，待添加)

默认启动端口`80`, 默认存储路径`/data/tmp/`

启动方式

```
go run main.go
```

## 功能
- [x] 单机文件存储
- [x] MySQL主从复制
- [x] 用户权限验证
- [x] 秒传功能
- [] 分块上传
- [] 断点续传


## 环境参数
Go `go version go1.10.2 linux/amd64`
 
操作系统 `Ubuntu 16.04`

IDE `VS code`
