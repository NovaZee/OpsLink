package informer

import (
	"fmt"
	"github.com/oppslink/protocol/logger"
	corev1 "k8s.io/api/core/v1"
	"sync"
)

type ServiceInformer struct {
	localCache sync.Map
}

func (s *ServiceInformer) OnAdd(obj interface{}, isInInitialList bool) {
	if svc, ok := obj.(*corev1.Service); ok {
		s.Add(svc)
	}
	logger.Debugw("service informer webhook", "action", "OnAdd", "Name", obj.(*corev1.Service).Name, "namespace", obj.(*corev1.Service).Namespace)
}

func (s *ServiceInformer) OnUpdate(oldObj, newObj interface{}) {
	err := s.Update(newObj.(*corev1.Service))
	if err != nil {
		logger.Warnw("service informer webhook", err, "action", "OnUpdate")
	}
	logger.Debugw("service informer webhook", "action", "OnUpdate", "Name", newObj.(*corev1.Service).Name, "namespace", newObj.(*corev1.Service).Namespace)
}

func (s *ServiceInformer) OnDelete(obj interface{}) {
	if svc, ok := obj.(*corev1.Service); ok {
		s.Delete(svc)
	}
	logger.Debugw("service informer webhook", "action", "OnDelete", "Name", obj.(*corev1.Service).Name, "namespace", obj.(*corev1.Service).Namespace)
}

func (s *ServiceInformer) Delete(service *corev1.Service) {
	if serviceList, ok := s.localCache.Load(service.Namespace); ok {
		list := serviceList.([]*corev1.Service)
		for k, needDeleteService := range list {
			if service.Name == needDeleteService.Name {
				newList := append(list[:k], list[k+1:]...)
				s.localCache.Store(service.Namespace, newList)
				break
			}
		}
	}
}

func (s *ServiceInformer) Add(service *corev1.Service) {
	if serviceList, ok := s.localCache.Load(service.Namespace); ok {
		serviceList = append(serviceList.([]*corev1.Service), service)
		s.localCache.Store(service.Namespace, serviceList)
	} else {
		newServiceList := make([]*corev1.Service, 0)
		newServiceList = append(newServiceList, service)
		s.localCache.Store(service.Namespace, newServiceList)
	}
}

func (s *ServiceInformer) Update(service *corev1.Service) error {
	if serviceList, ok := s.localCache.Load(service.Namespace); ok {
		list := serviceList.([]*corev1.Service)
		for k, needUpdateService := range list {
			if service.Name == needUpdateService.Name {
				list[k] = service
			}
		}
		return nil
	}

	return fmt.Errorf("service-%s update error", service.Name)
}

func (s *ServiceInformer) ListAll(namespace string) ([]*corev1.Service, error) {
	if serviceList, ok := s.localCache.Load(namespace); ok {
		return serviceList.([]*corev1.Service), nil
	}

	return []*corev1.Service{}, nil
}
