package handler

import (
	"cloudstore/config"
	"cloudstore/util"
	"context"
	"log"
	"net/http"

	userProto "cloudstore/service/account/proto"
	uploadProto "cloudstore/service/upload/proto"

	"github.com/gin-gonic/gin"
	micro "github.com/micro/go-micro"
)

var (
	userCli   userProto.UserService
	uploadCli uploadProto.UploadService
)

func init() {
	service := micro.NewService()
	// 初始化， 解析命令行参数等
	service.Init()
	// 初始化一个account服务的客户端
	userCli = userProto.NewUserService("go.micro.service.user", service.Client())
	uploadCli = uploadProto.NewUploadService("go.micro.service.upload", service.Client())
}

// SignupHandler : 响应注册页面
func SignupHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signup.html")
}

// DoSignupHandler : 处理注册post请求
func DoSignupHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	passwd := c.Request.FormValue("password")

	resp, err := userCli.SignUp(context.TODO(), &userProto.ReqSignUp{
		Username: username,
		Password: passwd,
	})

	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": resp.Code,
		"msg":  resp.Message,
	})
}

// SigninHandler : 响应登录页面
func SigninHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signin.html")
}

// DoSigninHandler : 处理登录post请求
func DoSigninHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")

	resp, err := userCli.SignIn(context.TODO(), &userProto.ReqSignIn{
		Username: username,
		Password: password,
	})

	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	if resp.Code != util.StatusOK {
		c.JSON(http.StatusOK, gin.H{
			"msg":  resp.Message,
			"code": resp.Code,
		})
		return
	}
	// TODO : 添加UploadEntry和 DownloadEntry

	// 登录成功，返回用户信息
	res := util.RespMsg{
		Code: resp.Code,
		Msg:  resp.Message,
		Data: struct {
			Location      string
			Username      string
			Token         string
			UploadEntry   string
			DownloadEntry string
		}{
			Location:      "/static/view/home.html",
			Username:      username,
			Token:         resp.Token,
			UploadEntry:   config.UploadEntry,
			DownloadEntry: config.DownloadEntry,
		},
	}
	c.Data(http.StatusOK, "application/json", res.JSONByte())
}

// UserInfoHandler ： 查询用户信息
func UserInfoHandler(c *gin.Context) {
	// 1. 解析请求参数
	username := c.Request.FormValue("username")
	// 2. RPC传送请求
	resp, err := userCli.UserInfo(context.TODO(), &userProto.ReqUserInfo{
		Username: username,
	})

	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	// 3. 组装并且响应用户数据
	res := util.RespMsg{
		Code: resp.Code,
		Msg:  resp.Message,
		Data: gin.H{
			"Username":   username,
			"CreateTime": resp.CreateTime,
			// TODO: 完善其他字段信息
			// "LastEditTime": resp.LastEditTime,
		},
	}
	c.Data(http.StatusOK, "application/json", res.JSONByte())
}
