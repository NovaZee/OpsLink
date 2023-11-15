package service

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type SignalService struct {
	upgrader websocket.Upgrader
}

//写消息发送的pb

func (s *SignalService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// reject non websocket requests
	if !websocket.IsWebSocketUpgrade(r) {
		w.WriteHeader(404)
		return
	}
	//参数校验
	//http升级为websocket
	//建立连接
	//处理请求
	//处理响应下发
	//心跳

}

func (s *SignalService) ValidateConn(r *http.Request) {
	// token校验
	// 权限校验
}
