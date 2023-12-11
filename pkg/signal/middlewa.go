package signal

import (
	"github.com/denovo/permission/pkg/util"
	"github.com/oppslink/protocol/logger"
	"net/http"
	"time"
)

type AuthMiddleware struct{}

func (a *AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.URL != nil && r.URL.Path == "/signal/validate" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
	//1:log 推送日志
	//2:exec  进入容器内命令行操作
	//3:signal  根据用户订阅资源 进行相应资源变动的推送

	var authToken string
	authToken = r.FormValue(tokenParam)
	if authToken == "" {
		return
	}
	token, err := util.ParseToken(authToken)
	if err != nil {
		logger.Infow("websocket connecting", "result", "fail", "reason", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if token.Valid() != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
