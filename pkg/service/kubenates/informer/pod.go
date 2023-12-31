package informer

import (
	"fmt"
	"github.com/oppslink/protocol/logger"
	corev1 "k8s.io/api/core/v1"
	"sync"
)

type PodInformer struct {
	localCache sync.Map
}

// OnAdd add event informer 当有新的对象被创建时，将会调用这个函数
func (d *PodInformer) OnAdd(obj interface{}, isInInitialList bool) {
	if dep, ok := obj.(*corev1.Pod); ok {
		d.Add(dep)
	}
	logger.Debugw("pod informer webhook", "action", "OnDelete", "Name", obj.(*corev1.Pod).Name, "namespace", obj.(*corev1.Pod).Namespace)
}

// OnUpdate update event informer 当对象被修改时，将会调用这个函数。
// - resyncPeriod: if non-zero, will re-list this often (you will get OnUpdate
// calls, even if nothing changed). Otherwise, re-list will be delayed as
func (d *PodInformer) OnUpdate(oldObj, newObj interface{}) {
	err := d.Update(newObj.(*corev1.Pod))
	if err != nil {
		logger.Warnw("pod informer webhook", err, "action", "OnUpdate")
	}
	logger.Debugw("pod informer webhook", "action", "OnUpdate", "Name", newObj.(*corev1.Pod).Name, "namespace", newObj.(*corev1.Pod).Namespace)
}

// OnDelete delete event informer 当对象被删除时，将会调用这个函数。
func (d *PodInformer) OnDelete(obj interface{}) {
	if dep, ok := obj.(*corev1.Pod); ok {
		d.Delete(dep)
	}
	logger.Debugw("pod informer webhook", "action", "OnDelete", "Name", obj.(*corev1.Pod).Name, "namespace", obj.(*corev1.Pod).Namespace)
}

// Add informer to local cache
func (d *PodInformer) Add(pod *corev1.Pod) {
	if pods, ok := d.localCache.Load(pod.Namespace); ok {
		pods = append(pods.([]*corev1.Pod), pod)
		d.localCache.Store(pod.Namespace, pods)
	} else {
		newPods := make([]*corev1.Pod, 0)
		newPods = append(newPods, pod)
		d.localCache.Store(pod.Namespace, newPods)
	}

}

// Update informer to local cache
func (d *PodInformer) Update(pod *corev1.Pod) error {
	if pods, ok := d.localCache.Load(pod.Namespace); ok {
		cacheList := pods.([]*corev1.Pod)
		for k, needUpdatePod := range cacheList {
			if pod.Name == needUpdatePod.Name {
				cacheList[k] = pod
			}
		}
		return nil
	}

	return fmt.Errorf("deployment-%s update error", pod.Name)
}

// Delete informer to local cache
func (d *PodInformer) Delete(pod *corev1.Pod) {
	if pods, ok := d.localCache.Load(pod.Namespace); ok {
		cacheList := pods.([]*corev1.Pod)
		for k, deletePod := range cacheList {
			if pod.Name == deletePod.Name {
				newList := append(cacheList[:k], cacheList[k+1:]...)
				d.localCache.Store(pod.Namespace, newList)
				break
			}
		}
	}
}

// ListTargetALl 内存中读取自定ns的pods列表
func (d *PodInformer) ListTargetALl(ns string) ([]*corev1.Pod, error) {
	if pods, ok := d.localCache.Load(ns); ok {
		return pods.([]*corev1.Pod), nil
	}

	return []*corev1.Pod{}, nil
}

// ListALl 内存中读取pods列表
func (d *PodInformer) ListALl() ([]*corev1.Pod, error) {
	var ret []*corev1.Pod
	d.localCache.Range(func(key, value interface{}) bool {
		ret = append(ret, value.([]*corev1.Pod)...)
		return true
	})
	return ret, nil
}
