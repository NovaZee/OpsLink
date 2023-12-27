package kubenates

import (
	"github.com/denovo/permission/pkg/util"
	"github.com/gorilla/websocket"
	"github.com/oppslink/protocol/logger"
	"io"
	corev1 "k8s.io/api/core/v1"
	"log"
	"net/http"
	"path"
	"time"
)

var Up websocket.Upgrader

func init() {
	Up = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}

func (k *K8sClient) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !websocket.IsWebSocketUpgrade(r) {
		w.WriteHeader(404)
		return
	}
	token := r.Header.Get("Authorization")
	if token == "" {
		w.WriteHeader(401)
		return
	}
	// 解析token 时间
	parseToken, err := util.ParseToken(token)
	if err != nil {
		w.WriteHeader(401)
		return
	}
	err = parseToken.Valid()
	if err != nil {
		w.WriteHeader(401)
		return
	}
	// 升级连接为 WebSocket
	conn, err := Up.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// 获取查询参数 根据参数判断是连接到k8s获取日志还是进入容器内部执行命令
	// 获取URL路径的最后一个部分
	action := path.Base(r.URL.Path)
	defer func() {
		logger.Debugw("disconnect logs websocket connection", "closeTime", time.Now(), "message", err)
		_ = conn.Close()
	}()
	switch action {
	case "exec":
		//2：进入容器内部执行命令  路径：/:ns/:pname/:cname/
	case "logs":
		err = k.Logs(r, conn)
		if err != nil {
			logger.Warnw("logs", err)
			return
		}
	default:
		// 未知操作
		w.WriteHeader(400)
		return
	}
}

func (k *K8sClient) Exec() {

}

func (k *K8sClient) Logs(r *http.Request, conn *websocket.Conn) error {
	//1：连接到k8s 获取指定pod指定推送的日志  路径：/:ns/:pname/:cname/
	// 获取查询参数
	values := r.URL.Query()
	ns := values.Get("ns")
	podName := values.Get("podName")
	container := values.Get("container")
	var tailLine int64 = 100
	req := k.Clientset.CoreV1().Pods(ns).GetLogs(podName, &corev1.PodLogOptions{
		Follow: true, Container: container, TailLines: &tailLine})
	reader, err := req.Stream(r.Context())
	// stream err 或者 EOF 退出
	if err != nil || err == io.EOF {
		return err
	}
	defer reader.Close()
	for {
		buf := make([]byte, 1024)
		n, err := reader.Read(buf)
		if err != nil {
			break
		}
		if n > 0 {
			err = conn.WriteMessage(websocket.TextMessage, buf[0:n])
			if err != nil {
				return err
			}
		}
	}
	return nil
}
