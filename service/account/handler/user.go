package handler

import (
	"cloudstore/config"
	mydb "cloudstore/db"
	proto "cloudstore/service/account/proto"
	"cloudstore/util"
	"context"
	"fmt"
	"time"
)

// User : 实现UserServiceHandler接口的对象
type User struct{}

// IsTokenValid : TODO check if token valid
func IsTokenValid(token string) bool {
	if len(token) != 40 {
		return false
	}
	// TODO: 判断token的时效性，是否过期
	// TODO: 从数据库表tbl_user_token查询username对应的token信息
	// TODO: 对比两个token是否一致
	return true
}

// GenToken : generate token of username
func GenToken(username string) string {
	timestamp := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + timestamp + config.TokenSalt))
	token := tokenPrefix + timestamp[:8]
	fmt.Printf("username: %s Token: %s\n", username, token)
	return token
}

// SignUp : 用户注册
func (u *User) SignUp(ctx context.Context, req *proto.ReqSignUp, resp *proto.RespSignUp) error {
	// 1. 获取username and password
	username := req.Username
	passwd := req.Password
	// 2. 判断参数是否合法
	if len(username) < 3 || len(passwd) < 5 {
		resp.Code = util.StatusRegisterFailed
		resp.Message = "请求参数无效"
		return nil
	}
	// 3. 对密码进行加盐及取MD5值加密
	encPasswd := util.MD5([]byte(passwd + config.PwdSalt))

	// 4. 向file表存储记录
	suc := mydb.UserSignUp(username, encPasswd)
	if suc {
		resp.Code = util.StatusOK
		resp.Message = "注册成功"
		return nil

	} else {
		resp.Code = util.StatusRegisterFailed
		resp.Message = "注册失败"
		return nil
	}
}

// SignIn : 用户登录
func (u *User) SignIn(ctx context.Context, req *proto.ReqSignIn, resp *proto.RespSignIn) error {
	// 1. 获取账号密码
	username := req.Username
	passwd := req.Password
	// 2. 判断参数是否合法
	if len(username) < 3 || len(passwd) < 5 {
		resp.Code = util.StatusLoginFailed
		resp.Message = "请求参数无效"
		return nil
	}

	// 3. 对密码进行加盐及取MD5值加密
	encPasswd := util.MD5([]byte(passwd + config.PwdSalt))

	// 4. 检查用户名以及密码是否在db中
	suc := mydb.UserSignin(username, encPasswd)
	if !suc {
		resp.Code = util.StatusLoginFailed
		resp.Message = "用户名或密码错误"
		return nil
	}

	// 5. 生成token并且存到db中
	token := GenToken(username)
	suc = mydb.UpdateToken(username, token)
	if !suc {
		resp.Code = util.StatusLoginFailed
		resp.Message = "生成Token错误"
		return nil
	}

	// 6. 登陆成功
	resp.Code = util.StatusOK
	resp.Message = "登陆成功"
	resp.Token = token
	return nil
}

// UserInfo : 获取用户信息
func (u *User) UserInfo(ctx context.Context, req *proto.ReqUserInfo, resp *proto.RespUserInfo) error {
	// 1. 获取参数
	username := req.Username
	// 2. 检索相关用户的信息
	user, err := mydb.GetUserInfo(username)
	if err != nil {
		resp.Code = util.StatusServerError
		resp.Message = "查询用户信息失败"
		return nil
	}
	// 3. 返回用户信息
	resp.Code = util.StatusOK
	resp.Username = user.Username
	resp.CreateTime = user.CreateTime
	resp.LastEditTime = user.LastEditTime
	resp.Status = int32(user.Status)
	// TODO: 需增加接口支持完善用户信息(email/phone等)
	return nil
}
