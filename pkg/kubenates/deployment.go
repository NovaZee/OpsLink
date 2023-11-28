package kubenates

import (
	"fmt"
	"github.com/oppslink/protocol/logger"
	v1 "k8s.io/api/apps/v1"
	"sync"
)

type DeploymentHandler struct {
	localCache sync.Map
}

// OnAdd add event informer 当有新的对象被创建时，将会调用这个函数
func (d *DeploymentHandler) OnAdd(obj interface{}, isInInitialList bool) {
	//if dep, ok := obj.(*v1.Deployment); ok {
	//	d.Add(dep)
	//}
	//ws推送
}

// OnUpdate update event informer 当对象被修改时，将会调用这个函数。
// - resyncPeriod: if non-zero, will re-list this often (you will get OnUpdate
// calls, even if nothing changed). Otherwise, re-list will be delayed as
func (d *DeploymentHandler) OnUpdate(oldObj, newObj interface{}) {
	err := d.Update(newObj.(*v1.Deployment))
	if err != nil {
		logger.Warnw("deployment informer webhook", err, "action", "OnUpdate")
	}
}

// OnDelete delete event informer 当对象被删除时，将会调用这个函数。
func (d *DeploymentHandler) OnDelete(obj interface{}) {
	if dep, ok := obj.(*v1.Deployment); ok {
		d.Delete(dep)
	}
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

// Delete informer to local cache
func (d *DeploymentHandler) Delete(deployment *v1.Deployment) {

	if deploymentList, ok := d.localCache.Load(deployment.Namespace); ok {
		list := deploymentList.([]*v1.Deployment)
		for k, needDeleteDeployment := range list {
			if deployment.Name == needDeleteDeployment.Name {
				newList := append(list[:k], list[k+1:]...)
				d.localCache.Store(deployment.Namespace, newList)
				break
			}
		}
	}
}

// Add informer to local cache
func (d *DeploymentHandler) Add(deployment *v1.Deployment) {

	if deploymentList, ok := d.localCache.Load(deployment.Namespace); ok {
		deploymentList = append(deploymentList.([]*v1.Deployment), deployment)
		d.localCache.Store(deployment.Namespace, deploymentList)
	} else {
		newDeploymentList := make([]*v1.Deployment, 0)
		newDeploymentList = append(newDeploymentList, deployment)
		d.localCache.Store(deployment.Namespace, newDeploymentList)
	}

}

// Update informer to local cache
func (d *DeploymentHandler) Update(deployment *v1.Deployment) error {
	if deploymentList, ok := d.localCache.Load(deployment.Namespace); ok {
		list := deploymentList.([]*v1.Deployment)
		for k, needUpdateDeployment := range list {
			if deployment.Name == needUpdateDeployment.Name {
				list[k] = deployment
			}
		}
		return nil

	}

	return fmt.Errorf("deployment-%s update error", deployment.Name)
}
