package service

import (
	"github.com/denovo/permission/pkg/util"
	"github.com/denovo/permission/protoc/pb"
	"github.com/gorilla/websocket"
	"github.com/oppslink/protocol/logger"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type SignalService struct {
	upgrader websocket.Upgrader

	onlineMember map[int64]*role.Role

	sync sync.RWMutex
}

func NewSignalService() *SignalService {
	return &SignalService{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		onlineMember: make(map[int64]*role.Role, 100),
	}
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
	defer func() {
		conn.Close()
		sLogger.Infow("Ending WS connection", "closeTime", time.Now())
	}()
	//参数校验
	//http升级为websocket
	//建立连接
	//处理请求
	//处理响应下发
	//心跳

}

func (s *SignalService) connected(role *role.Role) {
	s.sync.Lock()
	defer s.sync.Unlock()
	s.onlineMember[role.Id] = role
}

func (s *SignalService) validateConn(r *http.Request) (int64, string, error) {
	token := r.FormValue(tokenParam)
	parseToken, err := util.ParseToken(token)
	if err != nil {
		return 0, "", err
	}
	i, err := strconv.ParseInt(parseToken.Id, 0, 64)
	if err != nil {
		return 0, "", err
	}
	return i, parseToken.UserName, nil
}
