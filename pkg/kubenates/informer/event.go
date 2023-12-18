package informer

import (
	"fmt"
	"github.com/oppslink/protocol/logger"
	corev1 "k8s.io/api/core/v1"
	"sync"
)

type EventInformer struct {
	localCache sync.Map
}

func (e *EventInformer) OnAdd(obj interface{}, isInInitialList bool) {
	e.store(obj, false)
	logger.Infow("event informer webhook", "action", "OnAdd", "Message", obj.(*corev1.Event).Message, "namespace", obj.(*corev1.Event).Namespace)
}
func (e *EventInformer) OnUpdate(oldObj, newObj interface{}) {
	e.store(newObj, false)
	logger.Infow("event informer webhook", "action", "OnAdd", "Message", newObj.(*corev1.Event).Message, "namespace", newObj.(*corev1.Event).Namespace)
}
func (e *EventInformer) OnDelete(obj interface{}) {
	e.store(obj, true)
	logger.Infow("event informer webhook", "action", "OnAdd", "Message", obj.(*corev1.Event).Message, "namespace", obj.(*corev1.Event).Namespace)
}

func (e *EventInformer) store(obj interface{}, isdelete bool) {
	if event, ok := obj.(*corev1.Event); ok {
		key := fmt.Sprintf("%s_%s_%s", event.Namespace, event.InvolvedObject.Kind, event.InvolvedObject.Name)
		if !isdelete {
			e.localCache.Store(key, event)
		} else {
			e.localCache.Delete(key)
		}
	}
}

func (e *EventInformer) GetEvent(ns string, kind string, name string) string {
	key := fmt.Sprintf("%s_%s_%s", ns, kind, name)
	if v, ok := e.localCache.Load(key); ok {
		return v.(*corev1.Event).Message
	}
	return ""
}
