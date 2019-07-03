package rpc

import (
	"cloudstore/config"
	upProto "cloudstore/service/upload/proto"
	"context"
)

// Upload : upload结构体
type Upload struct{}

// UploadEntry : 获取上传入口
func (u *Upload) UploadEntry(ctx context.Context, req *upProto.ReqEntry, resp *upProto.RespEntry) error {
	resp.Entry = config.UploadEntry
	return nil
}
