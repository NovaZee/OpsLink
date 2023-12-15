package informer

import (
	"fmt"
	"github.com/oppslink/protocol/logger"
	v1 "k8s.io/api/core/v1"
	"sync"
)

type ConfigMapInformer struct {
	localCache sync.Map
}

// OnAdd add event informer 当有新的对象被创建时，将会调用这个函数
func (c *ConfigMapInformer) OnAdd(obj interface{}, isInInitialList bool) {
	if cm, ok := obj.(*v1.ConfigMap); ok {
		c.Add(cm)
	}
	logger.Infow("configmap informer webhook", "action", "OnAdd", "Name", obj.(*v1.ConfigMap).Name)
	//ws推送
}

// OnUpdate update event informer 当对象被修改时，将会调用这个函数。
func (c *ConfigMapInformer) OnUpdate(oldObj, newObj interface{}) {
	err := c.Update(newObj.(*v1.ConfigMap))
	if err != nil {
		logger.Warnw("configmap informer webhook", err, "action", "OnUpdate")
	}
	logger.Infow("configmap informer webhook", "action", "OnUpdate", "Name", newObj.(*v1.ConfigMap).Name)
}

// OnDelete delete event informer 当对象被删除时，将会调用这个函数。
func (c *ConfigMapInformer) OnDelete(obj interface{}) {
	if cm, ok := obj.(*v1.ConfigMap); ok {
		c.Delete(cm)
	}
	logger.Infow("configmap informer webhook", "action", "OnDelete", "Name", obj.(*v1.ConfigMap).Name)
}

// Delete informer to local cache
func (c *ConfigMapInformer) Delete(configmap *v1.ConfigMap) {
	if configmapList, ok := c.localCache.Load(configmap.Namespace); ok {
		list := configmapList.([]*v1.ConfigMap)
		for k, needDeleteConfigMap := range list {
			if configmap.Name == needDeleteConfigMap.Name {
				newList := append(list[:k], list[k+1:]...)
				c.localCache.Store(configmap.Namespace, newList)
				break
			}
		}
	}
}

// Add informer to local cache
func (c *ConfigMapInformer) Add(configmap *v1.ConfigMap) {
	if configmapList, ok := c.localCache.Load(configmap.Namespace); ok {
		configmapList = append(configmapList.([]*v1.ConfigMap), configmap)
		c.localCache.Store(configmap.Namespace, configmapList)
	} else {
		newConfigMapList := make([]*v1.ConfigMap, 0)
		newConfigMapList = append(newConfigMapList, configmap)
		c.localCache.Store(configmap.Namespace, newConfigMapList)
	}
}

// Update informer to local cache
func (c *ConfigMapInformer) Update(configmap *v1.ConfigMap) error {
	if configmapList, ok := c.localCache.Load(configmap.Namespace); ok {
		list := configmapList.([]*v1.ConfigMap)
		for k, needUpdateConfigMap := range list {
			if configmap.Name == needUpdateConfigMap.Name {
				list[k] = configmap
			}
		}
		return nil
	}

	return fmt.Errorf("configmap-%s update error", configmap.Name)
}

func (c *ConfigMapInformer) List(namespace string) ([]*v1.ConfigMap, error) {
	if configmapList, ok := c.localCache.Load(namespace); ok {
		return configmapList.([]*v1.ConfigMap), nil
	}

	return []*v1.ConfigMap{}, nil
}

func (c *ConfigMapInformer) Get(namespace, name string) (*v1.ConfigMap, error) {
	list, ok := c.localCache.Load(namespace)
	if !ok {
		return nil, fmt.Errorf("configmap-%s get error: not found in the cache", name)
	}
	lists := list.([]*v1.ConfigMap)
	for _, configMap := range lists {
		if configMap.Name == name {
			return configMap, nil
		}
	}

	return nil, nil
}
