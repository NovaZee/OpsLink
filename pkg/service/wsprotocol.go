package service

import (
	"github.com/denovo/permission/protoc/model"
	"github.com/denovo/permission/protoc/signal"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/oppslink/protocol/logger"
	"google.golang.org/protobuf/encoding/protojson"
	"sync"
	"time"
)

const (
	pingFrequency = 10 * time.Second
	pingTimeout   = 2 * time.Second
)
const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 10 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type WsSignalConnClient struct {
	hub *HubSet

	conn *websocket.Conn

	// Buffered channel of outbound messages.
	ResponseBuffer chan []byte

	transportType int

	role *model.Role

	mu sync.Mutex
}

func NewWsSignalConn(hub *HubSet, conn *websocket.Conn, id int64, name string) *WsSignalConnClient {
	return &WsSignalConnClient{
		hub:            hub,
		conn:           conn,
		mu:             sync.Mutex{},
		ResponseBuffer: make(chan []byte, 256),
		transportType:  0,
		role: &model.Role{
			Id:   id,
			Name: name,
		},
	}
}

func (wsc *WsSignalConnClient) readRequest() (*signal.SignalRequest, int, error) {
	for {
		// handle special messages and pass on the rest
		messageType, payload, err := wsc.conn.ReadMessage()
		if err != nil {
			return nil, 0, err
		}

		msg := &signal.SignalRequest{}
		switch messageType {
		case websocket.BinaryMessage:
			if wsc.transportType == 1 {
				wsc.mu.Lock()
				// switch to protobuf if client supports it
				wsc.transportType = 0
				wsc.mu.Unlock()
			}
			// protobuf encoded
			err := proto.Unmarshal(payload, msg)
			return msg, len(payload), err
		case websocket.TextMessage:
			wsc.mu.Lock()
			// json encoded, also write back JSON
			wsc.transportType = 1
			wsc.mu.Unlock()
			err := protojson.Unmarshal(payload, msg)
			return msg, len(payload), err
		default:
			logger.Debugw("unsupported message", "message", messageType)
			return nil, len(payload), nil
		}
	}
}

func (wsc *WsSignalConnClient) writeResponse(response *signal.SignalResponse) (int, error) {
	var msgType int
	var payload []byte
	var err error

	wsc.mu.Lock()
	defer wsc.mu.Unlock()

	if wsc.transportType == 1 {
		msgType = websocket.TextMessage
		payload, err = protojson.Marshal(response)
	} else {
		msgType = websocket.BinaryMessage
		payload, err = proto.Marshal(response)
	}
	if err != nil {
		return 0, err
	}

	return len(payload), wsc.conn.WriteMessage(msgType, payload)
}

func (wsc *WsSignalConnClient) pingWorker() {
	for {
		<-time.After(pingPeriod)
		err := wsc.conn.WriteControl(websocket.PingMessage, []byte(""), time.Now().Add(pingTimeout))
		if err != nil {
			return
		}
	}
}
