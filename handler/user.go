package handler

import (
	mydb "cloudstore/db"
	"cloudstore/util"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	pwdSalt   = "!(@*#&$^%"
	tokenSalt = "_tokenSalt"
)

// SignUpHandler : 注册界面
func SignUpHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signup.html")
}

// DoSignUpHandler : 注册用户
func DoSignUpHandler(c *gin.Context) {
	// 1. 获取用户名和密码
	username := c.Request.FormValue("username")
	passwd := c.Request.FormValue("password")
	// 2. 判断参数是否合法
	if len(username) < 3 || len(passwd) < 5 {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "请求参数无效",
			"code": util.StatusRegisterFailed,
		})
		return
	}
	// 3. 对密码进行加盐及取MD5值加密
	encPasswd := util.MD5([]byte(passwd + pwdSalt))

	// 4. 向file表存储记录
	suc := mydb.UserSignUp(username, encPasswd)
	if suc {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "注册成功",
			"code": util.StatusOK,
		})

	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "注册失败",
			"code": util.StatusRegisterFailed,
		})
	}
	return
}

// SignInHandler : 登陆页面
func SignInHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signin.html")
	return
}

// DoSignInHandler : sign in handler
func DoSignInHandler(c *gin.Context) {
	// 1. 获取账号密码
	username := c.Request.FormValue("username")
	passwd := c.Request.FormValue("password")
	// 2. 判断参数是否合法
	if len(username) < 3 || len(passwd) < 5 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "请求参数无效",
			"code": util.StatusLoginFailed,
		})
		return
	}

	// 3. 对密码进行加盐及取MD5值加密
	encPasswd := util.MD5([]byte(passwd + pwdSalt))

	// 4. 检查用户名以及密码是否在db中
	suc := mydb.UserSignin(username, encPasswd)
	if !suc {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "用户名或密码错误",
			"code": util.StatusLoginFailed,
		})
		return
	}

	// 5. 生成token并且存到db中
	token := GenToken(username)
	suc = mydb.UpdateToken(username, token)
	if !suc {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "生成Token错误",
			"code": util.StatusLoginFailed,
		})
		return
	}

	// 6. 登陆成功
	resp := util.NewRespMsg(util.StatusOK, "OK", struct {
		Location string
		Username string
		Token    string
	}{
		Location: "/static/view/home.html",
		Username: username,
		Token:    token,
	})
	fmt.Printf("user :%s, resp:%s\n", username, resp.JSONString())
	c.Data(http.StatusOK, "application/json", resp.JSONByte())
}

// UserInfoHandler : user info search handler
func UserInfoHandler(c *gin.Context) {
	// 1. 获取参数
	username := c.Request.FormValue("username")
	token := c.Request.FormValue("token")
	// 2. 判断token是否合法
	isValidToken := IsTokenValid(token)
	if !isValidToken {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "token无效，请重新登陆",
			"code": util.StatusInvalidToken,
		})
		return
	}
	// 3. 检索相关用户的信息
	user, err := mydb.GetUserInfo(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "查询用户信息失败",
			"code": util.StatusServerError,
		})
		return
	}
	// 4. 返回用户信息
	resp := util.NewRespMsg(util.StatusOK, "OK", user)
	c.Data(http.StatusOK, "application/json", resp.JSONByte())
}

// GenToken : generate token of username
func GenToken(username string) string {
	timestamp := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + timestamp + tokenSalt))
	token := tokenPrefix + timestamp[:8]
	fmt.Printf("username: %s Token: %s\n", username, token)
	return token
}

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
