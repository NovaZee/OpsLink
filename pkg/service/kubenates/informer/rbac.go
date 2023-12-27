package informer

import (
	"fmt"
	"github.com/oppslink/protocol/logger"
	rbacv1 "k8s.io/api/rbac/v1"
	"sort"
	"sync"
)

type RBACInformer struct {
	localCache sync.Map
}

// OnAdd add event informer
func (ri *RBACInformer) OnAdd(obj interface{}, isInInitialList bool) {
	if role, ok := obj.(*rbacv1.Role); ok {
		ri.Add(role)
	}
	logger.Debugw("rbac informer webhook", "action", "OnAdd", "Name", obj.(*rbacv1.Role).Name, "namespace", obj.(*rbacv1.Role).Namespace)
}

// OnUpdate update event informer
func (ri *RBACInformer) OnUpdate(oldObj, newObj interface{}) {
	err := ri.Update(newObj.(*rbacv1.Role))
	if err != nil {
		logger.Warnw("rbac informer webhook", err, "action", "OnUpdate")
	}
	logger.Debugw("rbac informer webhook", "action", "OnUpdate", "Name", newObj.(*rbacv1.Role).Name, "namespace", newObj.(*rbacv1.Role).Namespace)
}

// OnDelete delete event informer
func (ri *RBACInformer) OnDelete(obj interface{}) {
	if role, ok := obj.(*rbacv1.Role); ok {
		ri.Delete(role)
	}
	logger.Debugw("rbac informer webhook", "action", "OnDelete", "Name", obj.(*rbacv1.Role).Name, "namespace", obj.(*rbacv1.Role).Namespace)
}

// Add informer to local cache
func (ri *RBACInformer) Add(role *rbacv1.Role) {
	if roles, ok := ri.localCache.Load(role.Namespace); ok {
		roles = append(roles.([]*rbacv1.Role), role)
		ri.localCache.Store(role.Namespace, roles)
	} else {
		newRoles := make([]*rbacv1.Role, 0)
		newRoles = append(newRoles, role)
		ri.localCache.Store(role.Namespace, newRoles)
	}
}

// Update informer to local cache
func (ri *RBACInformer) Update(role *rbacv1.Role) error {
	if roles, ok := ri.localCache.Load(role.Namespace); ok {
		cacheList := roles.([]*rbacv1.Role)
		for k, needUpdateRole := range cacheList {
			if role.Name == needUpdateRole.Name {
				cacheList[k] = role
			}
		}
		return nil
	}

	return fmt.Errorf("role-%s update error", role.Name)
}

// Delete informer to local cache
func (ri *RBACInformer) Delete(role *rbacv1.Role) {
	if roles, ok := ri.localCache.Load(role.Namespace); ok {
		cacheList := roles.([]*rbacv1.Role)
		for k, deleteRole := range cacheList {
			if role.Name == deleteRole.Name {
				newList := append(cacheList[:k], cacheList[k+1:]...)
				ri.localCache.Store(role.Namespace, newList)
				break
			}
		}
	}
}

// ListTargetAll reads the list of roles in a specific namespace from memory
func (ri *RBACInformer) ListTargetAll(ns string) ([]*rbacv1.Role, error) {
	if roles, ok := ri.localCache.Load(ns); ok {
		return roles.([]*rbacv1.Role), nil
	}

	return []*rbacv1.Role{}, nil
}

// ListAll reads the list of all roles from memory
func (ri *RBACInformer) ListAll() ([]*rbacv1.Role, error) {
	var ret []*rbacv1.Role
	ri.localCache.Range(func(key, value interface{}) bool {
		ret = append(ret, value.([]*rbacv1.Role)...)
		return true
	})
	return ret, nil
}

func (ri *RBACInformer) ListAllByNs(ns string) []*rbacv1.Role {
	if list, ok := ri.localCache.Load(ns); ok {
		newList := list.([]*rbacv1.Role)
		sort.Sort(RbacV1(newList))
		return newList
	}

	return []*rbacv1.Role{}
}

type RbacV1 []*rbacv1.Role

func (r RbacV1) Len() int {
	return len(r)
}
func (r RbacV1) Less(i, j int) bool {
	return r[i].CreationTimestamp.Time.After(r[j].CreationTimestamp.Time)
}
func (r RbacV1) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
