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
	pingTimeout = 2 * time.Second
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 10 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type WsSignalConnClient struct {
	hub *HubSet

	conn *websocket.Conn

	// Buffered channel of outbound messages.
	ResponseBuffer chan []byte

	transportType int

	role *model.Role

	lastSignalTime int64

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
		lastSignalTime: time.Now().Unix(),
	}
}

// ReadRequest reads data from the connection.
// It handles special messages and passes on the rest.
// It determines the message type (protobuf or JSON) based on the payload.
// Returns a signal.SignalRequest message or an error.
func (wsc *WsSignalConnClient) ReadRequest() (*signal.SignalRequest, error) {
	for {
		// handle special messages and pass on the rest
		messageType, payload, err := wsc.conn.ReadMessage()
		if err != nil {
			return nil, err
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
			return msg, err
		case websocket.TextMessage:
			wsc.mu.Lock()
			// json encoded, also write back JSON
			wsc.transportType = 1
			wsc.mu.Unlock()
			err := protojson.Unmarshal(payload, msg)
			return msg, err
		default:
			logger.Debugw("unsupported message", "message", messageType)
			return nil, nil
		}
	}
}

// WriteResponse writes the provided signal.SignalResponse to the connection.
// The function determines the message type (protobuf or JSON) based on the transportType.
// Returns the number of bytes written and any encountered error.
func (wsc *WsSignalConnClient) WriteResponse(response *signal.SignalResponse) (int, error) {
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

// UpdateLastSignalTime update time for last ping.
func (wsc *WsSignalConnClient) UpdateLastSignalTime(signalTime int64) {
	wsc.lastSignalTime = signalTime
}

// pingWorker
func (wsc *WsSignalConnClient) pingWorker() {
	for {
		<-time.After(pingPeriod)
		err := wsc.conn.WriteControl(websocket.PingMessage, []byte(""), time.Now().Add(pingTimeout))
		if err != nil {
			return
		}
	}
}
