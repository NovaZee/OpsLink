package kubenates

import "sync"

type PodHandler struct {
	podData sync.Map
}

// OnAdd add event informer 当有新的对象被创建时，将会调用这个函数
func (d *PodHandler) OnAdd(obj interface{}, isInInitialList bool) {

}

// OnUpdate update event informer 当对象被修改时，将会调用这个函数。
// - resyncPeriod: if non-zero, will re-list this often (you will get OnUpdate
// calls, even if nothing changed). Otherwise, re-list will be delayed as
func (d *PodHandler) OnUpdate(oldObj, newObj interface{}) {

}

// OnDelete delete event informer 当对象被删除时，将会调用这个函数。
func (d *PodHandler) OnDelete(obj interface{}) {

}

// HandleAddEvent handler active add event
func (d *PodHandler) HandleAddEvent(data string) {

}

// HandleUpdateEvent handler active update event
func (d *PodHandler) HandleUpdateEvent(data string) {

}

// HandleDeleteEvent handler active delete event
func (d *PodHandler) HandleDeleteEvent(data string) {

}
