package informer

import (
	"sync"
)

type PodInformer struct {
	localCache sync.Map
}

// OnAdd add event informer 当有新的对象被创建时，将会调用这个函数
func (d *PodInformer) OnAdd(obj interface{}, isInInitialList bool) {

}

// OnUpdate update event informer 当对象被修改时，将会调用这个函数。
// - resyncPeriod: if non-zero, will re-list this often (you will get OnUpdate
// calls, even if nothing changed). Otherwise, re-list will be delayed as
func (d *PodInformer) OnUpdate(oldObj, newObj interface{}) {

}

// OnDelete delete event informer 当对象被删除时，将会调用这个函数。
func (d *PodInformer) OnDelete(obj interface{}) {

}
