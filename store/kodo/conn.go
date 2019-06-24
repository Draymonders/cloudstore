package kodo

import (
	"cloudstore/config"
	"context"
	"fmt"

	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

// PutObject : 上传数据
func PutObject(bucket string, localFile string, key string) bool {
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	// 秘钥对象
	mac := qbox.NewMac(config.KodoAccessKey, config.KodoSecretKey)
	// 生成上传token
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuadong
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false

	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	err := formUploader.PutFile(context.Background(), &ret, upToken, key, localFile, nil)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println(ret.Key, ret.Hash)
	return true
}

// GetObjectURL : 获取数据
func GetObjectURL(key string) string {
	publicAccessURL := storage.MakePublicURL(config.KodoDomain, key)
	return publicAccessURL
}
