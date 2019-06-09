package handler

import (
	"net/http"
	"util"
)

func HTTPInterceptor(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			username := r.Form.Get("username")
			token := r.Form.Get("token")

			// check user and token
			// if no error then next
			if len(username) < 3 || !IsTokenValid(token) {
				resp := util.NewRespMsg(
					util.StatusInvalidToken,
					"token无效",
					nil,
				)
				w.Write(resp.JSONByte())
				return
			}
			handler(w, r)
		})
}
