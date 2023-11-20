package service

import (
	"github.com/denovo/permission/pkg/util"
	role "github.com/denovo/permission/protoc/pb"
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

type SignalService struct {
	upgrader websocket.Upgrader

	onlineMember map[int64]*role.Role

	hub *HubSet

	sync sync.RWMutex
}

func NewSignalService() *SignalService {

	signalService := &SignalService{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		onlineMember: make(map[int64]*role.Role, 100),
		hub:          newHubSet(),
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
	id, name, _ := s.validateConn(r)
	// 升级连接为 WebSocket
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	//logger
	sLogger := LoggerWithRole(logger.GetLogger(), id, name)
	var role = &role.Role{
		Id:   id,
		Name: name,
	}
	s.connected(role)

	wsc := NewWsSignalConn(s.hub, conn)
	wsc.hub.register <- wsc

	//心跳
	go wsc.pingWorker()

	//处理请求
	go wsc.readRequest()
	//处理响应下发
	go wsc.writeResponse()

	sLogger.Infow("Ending WS connection", "closeTime", time.Now())

}

func (s *SignalService) connected(role *role.Role) {
	s.sync.Lock()
	defer s.sync.Unlock()
	s.onlineMember[role.Id] = role
}

func (s *SignalService) disConnected(roleId int64) {
	s.sync.Lock()
	defer s.sync.Unlock()
	delete(s.onlineMember, roleId)
}

func (s *SignalService) validateConn(r *http.Request) (int64, string, error) {
	token := r.FormValue(tokenParam)
	parseToken, err := util.ParseToken(token)
	if err != nil {
		return 0, "", err
	}
	if parseToken.UserID != 0 {
		return 0, "", errors.New("角色ID有误！")
	}

	return parseToken.UserID, parseToken.UserName, nil
}
