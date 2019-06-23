在这里，我使用了七牛云的对象存储Kodo，和阿里云的OSS,还有腾讯云的COS是同样的产品
## oss相关术语
[![Z9xLpn.md.jpg](https://s2.ax1x.com/2019/06/23/Z9xLpn.md.jpg)](https://imgchr.com/i/Z9xLpn)

## 包依赖关系解决

### unrecognized import path "golang.org/x/net/context" 解决方案
`$GOPATH`为项目的原路径

```bash
$ mkdir -p $GOPATH/src/golang.org/x/
$ cd $GOPATH/src/golang.org/x
$ git clone git clone https://github.com/golang/net.git net
$ go install net
```
## 开发文档
[七牛云开发文档](https://github.com/qiniu/api.v7/wiki)
