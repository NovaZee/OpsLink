package service

import (
	"github.com/denovo/permission/pkg/util"
	"net/http"
	"time"
)

const (
	pingFrequency = 10 * time.Second
	pingTimeout   = 2 * time.Second
)
const (
	tokenParam = "token"
)

type AuthMiddleware struct{}

func (a *AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.URL != nil && r.URL.Path == "/rtc/validate" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	var authToken string
	authToken = r.FormValue(tokenParam)
	if authToken == "" {
		return
	}
	token, err := util.ParseToken(tokenParam)
	if err != nil {
		return
	}
	token.
		//处理token
		//处理权限
		//加入对应群组
		next.ServeHTTP(w, r)
}
