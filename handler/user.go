package handler

import (
	mydb "cloudstore/db"
	"cloudstore/util"
	"fmt"
	"net/http"
	"time"
)

const (
	pwdSalt   = "!(@*#&$^%"
	tokenSalt = "_tokenSalt"
)

// SignUpHandler : sign up handler
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.Redirect(w, r, "/static/view/signup.html", http.StatusFound)
		return
	} else if r.Method == http.MethodPost {
		// parse form from post
		r.ParseForm()

		// 1. get username and password
		username := r.Form.Get("username")
		passwd := r.Form.Get("password")
		fmt.Printf("sign up user: %s passwd : %s\n", username, passwd)
		// check username and password valid
		if len(username) < 3 || len(passwd) < 5 {
			w.Write([]byte("账户密码设置不合法"))
			return
		}
		// 2. encry password user Md5 and salt
		encPasswd := util.MD5([]byte(passwd + pwdSalt))
		fmt.Printf("sign up user: %s encPasswd : %s\n", username, encPasswd)
		// 3. store username and password to DB
		suc := mydb.UserSignUp(username, encPasswd)
		if suc {
			w.Write([]byte("OK"))
		} else {
			w.Write([]byte("Sign Up error"))
		}
	}
}

// SignInHandler : sign in handler
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.Redirect(w, r, "/static/view/signin.html", http.StatusFound)
		return
	} else if r.Method == http.MethodPost {
		// parse form from post
		r.ParseForm()
		// 1. get username and password
		username := r.Form.Get("username")
		passwd := r.Form.Get("password")
		// check username and password valid
		if len(username) < 3 || len(passwd) < 5 {
			w.Write([]byte("Invalid parameter"))
			return
		}

		// 2. encry password user Md5 and salt
		encPasswd := util.MD5([]byte(passwd + pwdSalt))

		// 3. check if username exists in DB
		suc := mydb.UserSignin(username, encPasswd)
		if !suc {
			w.Write([]byte("Username or Passwrod error"))
			return
		}

		// 4. generate token and store to DB
		token := GenToken(username)
		suc = mydb.UpdateToken(username, token)
		if !suc {
			w.Write([]byte("token update failed"))
			return
		}

		// 5. sign in ok
		// store token to client
		resp := util.NewRespMsg(util.StatusOK, "OK", struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://" + r.Host + "/static/view/home.html",
			Username: username,
			Token:    token,
		})
		fmt.Printf("user :%s, resp:%s\n", username, resp.JSONString())
		w.Write(resp.JSONByte())
	}
}

// UserInfoHandler : user info search handler
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// 1. parse form from post
	username := r.Form.Get("username")
	token := r.Form.Get("token")
	// 2. check if token valid
	isValidToken := IsTokenValid(token)
	if !isValidToken {
		w.Write([]byte("token has valid, please log in again"))
		return
	}
	// 3. search user info
	user, err := mydb.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	resp := util.NewRespMsg(util.StatusOK, "OK", user)
	w.Write(resp.JSONByte())
}

// GenToken : generate token of username
func GenToken(username string) string {
	timestamp := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + timestamp + tokenSalt))
	token := tokenPrefix + timestamp[:8]
	fmt.Printf("username: %s Token: %s\n", username, token)
	return token
}

// isTokenValid : TODO check if token valid
func IsTokenValid(token string) bool {
	if len(token) != 40 {
		return false
	}
	// TODO: 判断token的时效性，是否过期
	// TODO: 从数据库表tbl_user_token查询username对应的token信息
	// TODO: 对比两个token是否一致
	return true
}
