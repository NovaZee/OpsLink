package service

import (
	"log"
	"net/http"
	"time"
)

const (
	pingFrequency = 10 * time.Second
	pingTimeout   = 2 * time.Second
)

func WSSignalConnection() {
}

type MuxHandler struct {
	handler Handler

	// next stores the next.ServeHTTP to reduce memory allocate
	next func(rw http.ResponseWriter, r *http.Request)
}

type Handler interface {
	ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

func (m MuxHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	m.handler.ServeHTTP(rw, r, m.next)
}

type AuthMiddleware struct{}

func (a *AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// 升级连接为 WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("Received message: %s\n", message)

		// 在这里处理收到的 WebSocket 消息
		// ...

		// 发送响应消息
		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Println(err)
			return
		}
	}
	//参数校验
	//http升级为websocket
	//建立连接
	//处理请求
	//处理响应下发
	//心跳
	next.ServeHTTP(w, r)
}
