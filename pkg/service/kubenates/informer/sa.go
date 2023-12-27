package informer

import (
	"fmt"
	"github.com/oppslink/protocol/logger"
	corev1 "k8s.io/api/core/v1"
	"sync"
)

type SAInformer struct {
	localCache sync.Map
}

// OnAdd add event informer
func (si *SAInformer) OnAdd(obj interface{}, isInInitialList bool) {
	if sa, ok := obj.(*corev1.ServiceAccount); ok {
		si.Add(sa)
	}
	logger.Debugw("sa informer webhook", "action", "OnAdd", "Name", obj.(*corev1.ServiceAccount).Name, "namespace", obj.(*corev1.ServiceAccount).Namespace)
}

// OnUpdate update event informer
func (si *SAInformer) OnUpdate(oldObj, newObj interface{}) {
	err := si.Update(newObj.(*corev1.ServiceAccount))
	if err != nil {
		logger.Warnw("sa informer webhook", err, "action", "OnUpdate")
	}
	logger.Debugw("sa informer webhook", "action", "OnUpdate", "Name", newObj.(*corev1.ServiceAccount).Name, "namespace", newObj.(*corev1.ServiceAccount).Namespace)
}

// OnDelete delete event informer
func (si *SAInformer) OnDelete(obj interface{}) {
	if sa, ok := obj.(*corev1.ServiceAccount); ok {
		si.Delete(sa)
	}
	logger.Debugw("sa informer webhook", "action", "OnDelete", "Name", obj.(*corev1.ServiceAccount).Name, "namespace", obj.(*corev1.ServiceAccount).Namespace)
}

// Add informer to local cache
func (si *SAInformer) Add(sa *corev1.ServiceAccount) {
	if sas, ok := si.localCache.Load(sa.Namespace); ok {
		sas = append(sas.([]*corev1.ServiceAccount), sa)
		si.localCache.Store(sa.Namespace, sas)
	} else {
		newSAs := make([]*corev1.ServiceAccount, 0)
		newSAs = append(newSAs, sa)
		si.localCache.Store(sa.Namespace, newSAs)
	}
}

// Update informer to local cache
func (si *SAInformer) Update(sa *corev1.ServiceAccount) error {
	if sas, ok := si.localCache.Load(sa.Namespace); ok {
		cacheList := sas.([]*corev1.ServiceAccount)
		for k, needUpdateSA := range cacheList {
			if sa.Name == needUpdateSA.Name {
				cacheList[k] = sa
			}
		}
		return nil
	}

	return fmt.Errorf("sa-%s update error", sa.Name)
}

// Delete informer to local cache
func (si *SAInformer) Delete(sa *corev1.ServiceAccount) {
	if sas, ok := si.localCache.Load(sa.Namespace); ok {
		cacheList := sas.([]*corev1.ServiceAccount)
		for k, deleteSA := range cacheList {
			if sa.Name == deleteSA.Name {
				newList := append(cacheList[:k], cacheList[k+1:]...)
				si.localCache.Store(sa.Namespace, newList)
				break
			}
		}
	}
}

// ListTargetAll reads the list of service accounts in a specific namespace from memory
func (si *SAInformer) ListTargetAll(ns string) []*corev1.ServiceAccount {
	if sas, ok := si.localCache.Load(ns); ok {
		return sas.([]*corev1.ServiceAccount)
	}

	return []*corev1.ServiceAccount{}
}

// ListAll reads the list of all service accounts from memory
func (si *SAInformer) ListAll() ([]*corev1.ServiceAccount, error) {
	var ret []*corev1.ServiceAccount
	si.localCache.Range(func(key, value interface{}) bool {
		ret = append(ret, value.([]*corev1.ServiceAccount)...)
		return true
	})
	return ret, nil
}
