package service

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type SignalService struct {
	upgrader websocket.Upgrader
}

func NewSignalService() *SignalService {
	return &SignalService{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

//写消息发送的pb

func (s *SignalService) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// 升级连接为 WebSocket
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
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
