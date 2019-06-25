package handler

import (
	"cloudstore/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HTTPInterceptor : 权限认证
func HTTPInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Request.FormValue("username")
		token := c.Request.FormValue("token")
		if len(username) < 3 || !IsTokenValid(token) {
			// 直接终止，不接着往下走了
			c.Abort()
			resp := util.NewRespMsg(
				util.StatusInvalidToken,
				"token无效",
				nil,
			)
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		c.Next()
	}
}
