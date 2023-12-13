package informer

import (
	"fmt"
	"github.com/oppslink/protocol/logger"
	corev1 "k8s.io/api/core/v1"
	"sync"
)

type NodeInformer struct {
	localCache sync.Map
}

func (n *NodeInformer) OnAdd(obj interface{}, isInInitialList bool) {
	if node, ok := obj.(*corev1.Node); ok {
		n.Add(node)
	}
	logger.Infow("node informer webhook", "action", "OnAdd", "Name", obj.(*corev1.Node).Name)
}

func (n *NodeInformer) OnUpdate(oldObj, newObj interface{}) {
	err := n.Update(newObj.(*corev1.Node))
	if err != nil {
		logger.Warnw("node informer webhook", err, "action", "OnUpdate")
	}
	logger.Infow("node informer webhook", "action", "OnUpdate", "Name", newObj.(*corev1.Node).Name)
}

func (n *NodeInformer) OnDelete(obj interface{}) {
	if node, ok := obj.(*corev1.Node); ok {
		n.Delete(node)
	}
	logger.Infow("node informer webhook", "action", "OnDelete", "Name", obj.(*corev1.Node).Name)
}

func (n *NodeInformer) Add(node *corev1.Node) {
	n.localCache.Store(node.Name, node)
}

func (n *NodeInformer) Update(node *corev1.Node) error {
	_, ok := n.localCache.Load(node.Name)
	if !ok {
		return fmt.Errorf("node-%s update error: not found in the cache", node.Name)
	}
	n.localCache.Store(node.Name, node)
	return nil
}

func (n *NodeInformer) Delete(node *corev1.Node) {
	n.localCache.Delete(node.Name)
}

func (n *NodeInformer) ListAll() []*corev1.Node {
	var nodes []*corev1.Node
	n.localCache.Range(func(key, value interface{}) bool {
		if node, ok := value.(*corev1.Node); ok {
			nodes = append(nodes, node)
		}
		return true
	})
	return nodes
}

// GetPodsNum 根据节点名称 获取pods数量
func (n *NodeInformer) GetPodsNum(node string) (num int) {
	n.localCache.Range(func(key, value interface{}) bool {
		list := value.([]*corev1.Pod)
		for _, pod := range list {
			if pod.Spec.NodeName == node {
				num++
			}
		}
		return true
	})
	return
}
