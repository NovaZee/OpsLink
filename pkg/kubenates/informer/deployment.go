package informer

import (
	"fmt"
	"github.com/oppslink/protocol/logger"
	v1 "k8s.io/api/apps/v1"
	"sync"
)

type DeploymentInformer struct {
	localCache sync.Map
}

// OnAdd add event informer 当有新的对象被创建时，将会调用这个函数
func (d *DeploymentInformer) OnAdd(obj interface{}, isInInitialList bool) {
	if dep, ok := obj.(*v1.Deployment); ok {
		d.Add(dep)
	}
	// 遍历 sync.Map 并打印所有项的键值对
	d.localCache.Range(func(key, value interface{}) bool {
		fmt.Println("Key:", key, "Value:", value)
		return true
	})
	//ws推送
}

// OnUpdate update event informer 当对象被修改时，将会调用这个函数。
// - resyncPeriod: if non-zero, will re-list this often (you will get OnUpdate
// calls, even if nothing changed). Otherwise, re-list will be delayed as
func (d *DeploymentInformer) OnUpdate(oldObj, newObj interface{}) {
	err := d.Update(newObj.(*v1.Deployment))
	if err != nil {
		logger.Warnw("deployment informer webhook", err, "action", "OnUpdate")
	}
}

// OnDelete delete event informer 当对象被删除时，将会调用这个函数。
func (d *DeploymentInformer) OnDelete(obj interface{}) {
	if dep, ok := obj.(*v1.Deployment); ok {
		d.Delete(dep)
	}
	// 遍历 sync.Map 并打印所有项的键值对
	d.localCache.Range(func(key, value interface{}) bool {
		fmt.Println("Key:", key, "Value:", value)
		return true
	})
	logger.Infow("deployment informer webhook", "action", "OnDelete", "Name", obj.(*v1.Deployment).Name)
}

// Delete informer to local cache
func (d *DeploymentInformer) Delete(deployment *v1.Deployment) {

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
func (d *DeploymentInformer) Add(deployment *v1.Deployment) {

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
func (d *DeploymentInformer) Update(deployment *v1.Deployment) error {
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

// ListALl 内存中读取deploymentList
func (d *DeploymentInformer) ListALl(namespace string) ([]*v1.Deployment, error) {
	if deploymentList, ok := d.localCache.Load(namespace); ok {
		return deploymentList.([]*v1.Deployment), nil
	}

	return []*v1.Deployment{}, nil
}
