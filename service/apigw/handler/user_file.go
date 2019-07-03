package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	userProto "cloudstore/service/account/proto"
)

// FilesQueryHandler : 文件查询，通过用户名和limit
func FilesQueryHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	limit, _ := strconv.Atoi(c.Request.FormValue("limit"))
	rpcResp, err := userCli.UserFiles(context.TODO(), &userProto.ReqUserFile{
		Username: username,
		Limit:    int32(limit),
	})

	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	if len(rpcResp.FileData) <= 0 {
		rpcResp.FileData = []byte("[]")
	}
	c.Data(http.StatusOK, "application/json", rpcResp.FileData)
}
