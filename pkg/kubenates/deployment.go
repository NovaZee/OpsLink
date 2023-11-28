package kubenates

import (
	"github.com/oppslink/protocol/logger"
	v1 "k8s.io/api/apps/v1"
	"sync"
)

type DeploymentHandler struct {
	deploymentData sync.Map
}

// OnAdd add event informer 当有新的对象被创建时，将会调用这个函数
func (d *DeploymentHandler) OnAdd(obj interface{}, isInInitialList bool) {

}

// OnUpdate update event informer 当对象被修改时，将会调用这个函数。
// - resyncPeriod: if non-zero, will re-list this often (you will get OnUpdate
// calls, even if nothing changed). Otherwise, re-list will be delayed as
func (d *DeploymentHandler) OnUpdate(oldObj, newObj interface{}) {

}

// OnDelete delete event informer 当对象被删除时，将会调用这个函数。
func (d *DeploymentHandler) OnDelete(obj interface{}) {
	logger.Infow("deployment informer webhook", "action", "OnDelete", "Name", obj.(*v1.Deployment).Name)
}

// HandleAddEvent handler active add event
func (d *DeploymentHandler) HandleAddEvent(data string) {

}

// HandleUpdateEvent handler active update event
func (d *DeploymentHandler) HandleUpdateEvent(data string) {

}

// HandleDeleteEvent handler active delete event
func (d *DeploymentHandler) HandleDeleteEvent(data string) {

}
