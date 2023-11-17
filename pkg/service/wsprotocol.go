package service

import (
	"github.com/denovo/permission/pkg/util"
	"github.com/oppslink/protocol/logger"
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
	if r.URL != nil && r.URL.Path == "/signal/validate" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	var authToken string
	authToken = r.FormValue(tokenParam)
	if authToken == "" {
		return
	}
	token, err := util.ParseToken(authToken)
	if err != nil {
		return
	}
	if token.Valid() != nil {
		return
	}

	next.ServeHTTP(w, r)
}

// LoggerWithRole logger util
func LoggerWithRole(l logger.Logger, id int64, name string) logger.Logger {
	values := make([]interface{}, 0, 4)
	if name != "" {
		values = append(values, "roleName", name)
	}
	if id != 0 {
		values = append(values, "roleId", id)
	}
	values = append(values, "connectionTime", time.Now())
	// enable sampling per participant
	return l.WithValues(values...)
}
