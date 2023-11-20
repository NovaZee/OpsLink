package service

import (
	"github.com/denovo/permission/pkg/util"
	"github.com/denovo/permission/protoc/signal"
	"github.com/gorilla/websocket"
	"github.com/oppslink/protocol/logger"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"sync"
	"time"
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

type SignalService struct {
	upgrader websocket.Upgrader

	// not support sent to target conn
	hub *HubSet

	sync sync.RWMutex
}

func NewSignalService() *SignalService {

	signalService := &SignalService{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		hub: newHubSet(),
	}
	//启动hub上下线协程
	go signalService.hub.run()

	return signalService
}

//写消息发送的pb

func (s *SignalService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// reject non websocket requests
	if !websocket.IsWebSocketUpgrade(r) {
		w.WriteHeader(404)
		return
	}
	token := r.FormValue(tokenParam)
	id, name, _ := s.validateConn(token)
	// 升级连接为 WebSocket
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	//logger
	sLogger := LoggerWithRole(logger.GetLogger(), id, name)

	wsc := NewWsSignalConn(s.hub, conn, id, name)
	wsc.hub.register <- wsc

	//心跳
	go wsc.pingWorker()
	// 只处理心跳 和 token续期
	//处理请求
	go func() {
		defer func() {
			wsc.hub.unregister <- wsc
			wsc.conn.Close()
		}()
		wsc.conn.SetReadLimit(maxMessageSize)
		wsc.conn.SetReadDeadline(time.Now().Add(pongWait))
		wsc.conn.SetPongHandler(func(string) error { wsc.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
		for {
			req, _, err := wsc.readRequest()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					logger.Warnw("error: %v", err)
				}
				break
			}
			switch m := req.Message.(type) {
			case *signal.SignalRequest_Ping:
				wsc.writeResponse(&signal.SignalResponse{
					Message: &signal.SignalResponse_Pong{
						Pong: time.Now().UnixMilli(),
					},
				})
			case *signal.SignalRequest_Renewal:
				refreshToken, err := util.RefreshToken(m.Renewal.Token)
				logger.Warnw("refreshToken error ", err)
				wsc.writeResponse(&signal.SignalResponse{
					Message: &signal.SignalResponse_RenewalResp{
						RenewalResp: &signal.RefreshToken{
							Token: refreshToken,
						},
					},
				})
			}
		}
	}()
	//处理响应下发
	go func() {
		//req, count, err := wsc.readRequest()
		defer func() {
			_ = wsc.conn.Close()
		}()
		for {
			select {
			case message, ok := <-wsc.ResponseBuffer:
				wsc.conn.SetWriteDeadline(time.Now().Add(writeWait))
				if !ok {
					// The hub closed the channel.
					wsc.conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}

				w, err := wsc.conn.NextWriter(websocket.TextMessage)
				if err != nil {
					return
				}
				w.Write(message)

				// Add queued chat messages to the current websocket message.
				n := len(wsc.ResponseBuffer)
				for i := 0; i < n; i++ {
					w.Write(newline)
					w.Write(<-wsc.ResponseBuffer)
				}

				if err := w.Close(); err != nil {
					return
				}
			}
		}
	}()

	sLogger.Infow("Ending WS connection", "closeTime", time.Now())

}

func (s *SignalService) validateConn(token string) (int64, string, error) {

	parseToken, err := util.ParseToken(token)
	if err != nil {
		return 0, "", err
	}
	if parseToken.UserID != 0 {
		return 0, "", errors.New("角色ID有误！")
	}

	return parseToken.UserID, parseToken.UserName, nil
}
