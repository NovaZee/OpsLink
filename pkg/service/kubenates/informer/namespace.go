package informer

import (
	"fmt"
	"github.com/oppslink/protocol/logger"
	corev1 "k8s.io/api/core/v1"
	"sync"
)

type NamespaceInformer struct {
	localCache sync.Map
}

// OnAdd add event informer 当有新的对象被创建时，将会调用这个函数
func (n *NamespaceInformer) OnAdd(obj interface{}, isInInitialList bool) {
	if ns, ok := obj.(*corev1.Namespace); ok {
		n.Add(ns)
	}
	logger.Infow("namespace informer webhook", "action", "OnAdd", "Name", obj.(*corev1.Namespace).Name)
}

// OnUpdate update event informer 当对象被修改时，将会调用这个函数。
func (n *NamespaceInformer) OnUpdate(oldObj, newObj interface{}) {
	err := n.Update(newObj.(*corev1.Namespace))
	if err != nil {
		logger.Warnw("namespace informer webhook", err, "action", "OnUpdate")
	}
	logger.Infow("namespace informer webhook", "action", "OnUpdate", "Name", newObj.(*corev1.Namespace).Name)
}

// OnDelete delete event informer 当对象被删除时，将会调用这个函数。
func (n *NamespaceInformer) OnDelete(obj interface{}) {
	if ns, ok := obj.(*corev1.Namespace); ok {
		n.Delete(ns)
	}
	logger.Infow("namespace informer webhook", "action", "OnDelete", "Name", obj.(*corev1.Namespace).Name)
}

// Add informer to local cache
func (n *NamespaceInformer) Add(namespace *corev1.Namespace) {
	n.localCache.Store(namespace.Name, namespace)
}

// Update informer to local cache
func (n *NamespaceInformer) Update(namespace *corev1.Namespace) error {
	_, ok := n.localCache.Load(namespace.Name)
	if !ok {
		return fmt.Errorf("namespace-%s update error: not found in the cache", namespace.Name)
	}
	n.localCache.Store(namespace.Name, namespace)
	return nil
}

// Delete informer from local cache
func (n *NamespaceInformer) Delete(namespace *corev1.Namespace) {
	n.localCache.Delete(namespace.Name)
}

// GetNamespace from local cache
func (n *NamespaceInformer) GetNamespace(name string) (*corev1.Namespace, error) {
	ns, ok := n.localCache.Load(name)
	if !ok {
		return nil, fmt.Errorf("namespace-%s get error: not found in the cache", name)
	}
	return ns.(*corev1.Namespace), nil
}

func (n *NamespaceInformer) ListAll(weight int) []*corev1.Namespace {
	var namespaces []*corev1.Namespace
	n.localCache.Range(func(key, value interface{}) bool {
		if ns, ok := value.(*corev1.Namespace); ok {
			// 权重为-1 list所有值，如果不是-1 kube-public和kube-system不对外开放
			if weight == -1 {
				namespaces = append(namespaces, ns)
			} else if ns.Name != "kube-public" && ns.Name != "kube-system" {
				namespaces = append(namespaces, ns)
			}
		}
		return true
	})
	return namespaces
}
