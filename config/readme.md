## ReadMe
想要启动这个Demo，首先要配置自己的以下信息。请照着 [doc目录](../doc/)进行运维部署

```
package config

// HOST主机信息
const (
	HOST = 你的IP地址
)

// MySQL
const (
	// MysqlLink : mysql 路由
	MysqlLink = "mysql账号:mysql密码@tcp(" + HOST + ":mysql端口)/icloud?charset=utf8"
)

// Redis
const (
	// RedisHost : redis 路由
	RedisHost = HOST + ":" + Redis端口
	// RedisPass : redis auth
	RedisPass = Redis密码
)

// Ceph
const (
	// CephAccessKey : 访问Key
	CephAccessKey = ""
	// CephSecretKey : 访问密钥
	CephSecretKey = ""
	// CephGWEndpoint : gateway地址
	CephGWEndpoint = "http://" + HOST + ":9080"
)

// kodo
const (
	// KodoBucket : bucket桶
	KodoBucket = ""
	// KodoEndpoint : kodo endpoint
	KodoEndpoint = ""
	// KodoAccesskey : kodo访问key
	KodoAccessKey = ""
	// KodoSecretKey : kodo访问key secret
	KodoSecretKey = ""
	// KodoDomain : kodo文件域名，可以自己绑定，也可以qiniu临时域名
	KodoDomain = ""
)

// RabbitMQ
const (
	// AsyncTransferEnable ： 异步上传文件
	AsyncTransferEnable = true
	// TransExchangeName : 用于文件transfer的交换机
	TransExchangeName = ""
	// TransKodoQueueName : oss转移队列名
	TransKodoQueueName = ""
	// TransKodoErrQueueName : oss转移失败后写入另一个队列的队列名
	TransKodoErrQueueName = ""
	// TransKodoRoutingKey : routingkey
	TransKodoRoutingKey = ""
	// RabbitURL : rabbitmq服务的入口url
	RabbitURL = "amqp://guest:guest@" + HOST + ":5672/"
)

// pwdSalt
const (
    // pwdSalt : 秘钥Salt
	pwdSalt   = ""
    // tokenSalt : TokenSalt
	tokenSalt = "_tokenSalt"
)
```