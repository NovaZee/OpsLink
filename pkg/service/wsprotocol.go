package service

import (
	"github.com/gorilla/websocket"
	"time"
)

const (
	pingFrequency = 10 * time.Second
	pingTimeout   = 2 * time.Second
)

type WsSignalConnClient struct {
	hub *HubSet

	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

func NewWsSignalConn(hub *HubSet, conn *websocket.Conn) *WsSignalConnClient {
	return &WsSignalConnClient{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}
}

func (wsc *WsSignalConnClient) readRequest() {

}

func (wsc *WsSignalConnClient) writeResponse() {

}

func (wsc *WsSignalConnClient) pingWorker() {
	for {
		<-time.After(pingFrequency)
		err := wsc.conn.WriteControl(websocket.PingMessage, []byte(""), time.Now().Add(pingTimeout))
		if err != nil {
			return
		}
	}
}
