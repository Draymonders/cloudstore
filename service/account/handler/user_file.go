package handler

import (
	mydb "cloudstore/db"
	proto "cloudstore/service/account/proto"
	"cloudstore/util"
	"context"
	"encoding/json"
	"fmt"
)

// DownloadURLHandler : generate the file url
func DownloadURLHandler(username, filehash string) string {
	tmpURL := fmt.Sprintf(
		"http://%s/file/download?hash=%s&username=%s", filehash, username)
	return tmpURL
}

// UserFiles : 获取用户文件列表
func (u *User) UserFiles(ctx context.Context, req *proto.ReqUserFile, resp *proto.RespUserFile) error {
	username := req.Username
	count := req.Limit
	userFiles, err := mydb.QueryUserFileMetas(username, (int)(count))
	if err != nil {
		resp.Code = util.StatusQueryUserFilesError
		resp.Message = "查询 user file表失败"
		return nil
	}

	data, _ := json.Marshal(userFiles)
	for i, ufile := range userFiles {
		userFiles[i].DownLoadUrl = DownloadURLHandler(username, ufile.Hash)
	}
	resp.FileData = data
	return nil
}
