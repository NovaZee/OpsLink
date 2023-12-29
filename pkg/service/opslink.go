package service

import (
	config "github.com/denovo/permission/config"
	"github.com/denovo/permission/pkg/service/casbin"
	opskube "github.com/denovo/permission/pkg/service/kubenates"
	opsstore "github.com/denovo/permission/pkg/service/store"
	"github.com/denovo/permission/pkg/signal"
	"github.com/gorilla/mux"
	"net/http"
)

type OpsLinkServer struct {
	Config        *config.OpsLinkConfig
	Casbin        *casbin.Casbin
	StoreService  opsstore.StoreService
	SignalService *signal.SignalService

	K8sClient *opskube.K8sClient

	doneChan   chan struct{}
	closedChan chan struct{}
}

func NewOpsLinkServer(config *config.OpsLinkConfig,
	casbin *casbin.Casbin,
	store opsstore.StoreService,
	kcs *opskube.K8sClient,
	ws *signal.SignalService,
) (os *OpsLinkServer, err error) {
	os = &OpsLinkServer{
		Config:        config,
		Casbin:        casbin,
		StoreService:  store,
		K8sClient:     kcs,
		SignalService: ws,
		closedChan:    make(chan struct{}),
	}

	r := mux.NewRouter()
	r.HandleFunc("/v1/logs", kcs.ServeHTTP)
	r.HandleFunc("/v1/exec", kcs.ServeHTTP)
	auth := &signal.AuthMiddleware{}
	m := &signal.MuxHandler{
		Handler: auth,
		Next:    ws.ServeHTTP,
	}
	r.HandleFunc("/signal/validate", m.ServeHTTP)
	http.Handle("/", r)
	go func() {
		http.ListenAndServe(":8085", r)
	}()
	return
}

func (server *OpsLinkServer) Start() error {
	server.doneChan = make(chan struct{})

	<-server.doneChan

	close(server.closedChan)
	return nil
}

func (server *OpsLinkServer) Stop(force bool) {
	//todo:如果使用本地内存启动，关闭之前等待数据同步
	//Before closing, check if there is any unsynchronized data.
	server.StoreService.Stop()
	// casbin 内存数据刷新csv
	server.Casbin.Enforcer.SavePolicy()
	close(server.doneChan)

	// wait for fully closed
	<-server.closedChan
}
