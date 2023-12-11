package kubeservice

import (
	"github.com/denovo/permission/pkg/kubenates/informer"
	"k8s.io/client-go/kubernetes"
)

type EventService struct {
	Ei     *informer.EventInformer
	Client kubernetes.Interface

	helper *helper
}

func NewEventService(client kubernetes.Interface) *EventService {
	return &EventService{Ei: &informer.EventInformer{}, Client: client, helper: &helper{}}
}

func (es *EventService) GetEvent(ns string, kind string, name string) (event string) {
	return es.Ei.GetEvent(ns, kind, name)
}
