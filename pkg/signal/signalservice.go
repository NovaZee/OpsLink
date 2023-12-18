package signal

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
	var err error
	token := r.FormValue(tokenParam)
	id, name, err := s.validateConn(token)
	if err != nil {
		logger.Debugw("validateConn", "error", err.Error())
	}
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
		wsc.conn.SetReadLimit(maxMessageSize)
		wsc.conn.SetReadDeadline(time.Now().Add(pongWait))
		wsc.conn.SetPongHandler(func(string) error { wsc.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
		for {
			req, errs := wsc.ReadRequest()
			if errs != nil {
				if websocket.IsUnexpectedCloseError(errs, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					sLogger.Warnw("error: %v", errs)
				}
				break
			}
			if req == nil {
				return
			}
			switch m := req.Message.(type) {
			case *signal.SignalRequest_Ping:
				wsc.UpdateLastSignalTime(m.Ping)
				wsc.WriteResponse(&signal.SignalResponse{
					Message: &signal.SignalResponse_Pong{
						Pong: time.Now().UnixMilli(),
					},
				})
			case *signal.SignalRequest_Renewal:
				refreshToken, err1 := s.renewal(id, m.Renewal.Token)
				if err1 != nil {
					sLogger.Infow("refreshToken ", "err", err1)
					return
				}
				sLogger.Debugw(" refreshToken ", "id", id)
				wsc.WriteResponse(&signal.SignalResponse{
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
		ticker := time.NewTicker(3)
		//req, count, err := wsc.readRequest()
		defer func() {
			sLogger.Infow("ending websocket connection", "closeTime", time.Now())
			wsc.hub.unregister <- wsc
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

				iow, err2 := wsc.conn.NextWriter(websocket.TextMessage)
				if err2 != nil {
					return
				}
				iow.Write(message)
				// Add queued chat messages to the current websocket message.
				n := len(wsc.ResponseBuffer)
				for i := 0; i < n; i++ {
					iow.Write(<-wsc.ResponseBuffer)
				}

				if err3 := iow.Close(); err != nil {
					sLogger.Infow("io.WriteCloser ", "err", err3)
					return
				}
			case <-ticker.C:
				if time.Now().Unix()-wsc.lastSignalTime > 5 {
					logger.Infow("websocket ping gap time out", "id", id)
					return
				}

			}
		}
	}()

}

func (s *SignalService) validateConn(token string) (int64, string, error) {

	parseToken, err := util.ParseToken(token)
	if err != nil {
		return 0, "", err
	}
	if parseToken.UserID == 0 {
		return 0, "", errors.New("角色ID有误！")
	}

	return parseToken.UserID, parseToken.UserName, nil
}

func (s *SignalService) renewal(id int64, refresh string) (string, error) {
	token, err := util.ParseToken(refresh)
	if err != nil {
		return "", err
	}
	if id == token.UserID {
		refreshToken, err := util.RefreshToken(refresh)
		if err != nil {
			return "", err
		}
		return refreshToken, nil
	} else {
		return "", errors.New("token not match！")
	}

}
