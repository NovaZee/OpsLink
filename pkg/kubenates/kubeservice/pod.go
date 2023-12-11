package kubeservice

import (
	"context"
	"github.com/denovo/permission/pkg/kubenates/informer"
	"github.com/denovo/permission/protoc/kube"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodService struct {
	Pi     *informer.PodInformer
	Client kubernetes.Interface

	EventHandler *EventService

	helper *helper
}

func NewPodService(client kubernetes.Interface, eh *EventService) *PodService {
	return &PodService{Pi: &informer.PodInformer{}, Client: client, helper: &helper{}, EventHandler: eh}
}

func (ps *PodService) GetDetail(ctx context.Context, ns, name string) (*v1.Pod, error) {
	get, err := ps.Client.CoreV1().Pods(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return get, nil
}

// GetByLabelInCache kubectl get pods -l app=nginx -n default
func (ps *PodService) GetByLabelInCache(ns, label string) (res []*kube.Pod, err error) {

	pods, err := ps.Pi.ListALl(ns)

	if err != nil {
		return
	}
	for _, pod := range pods {
		if pod.Labels["app"] == label {
			res = append(res, &kube.Pod{
				Name:         pod.Name,
				Namespace:    pod.Namespace,
				Images:       ps.helper.GetImagesByPod(pod.Spec.Containers),
				NodeName:     pod.Spec.NodeName,
				Phase:        string(pod.Status.Phase),
				Ip:           []string{pod.Status.PodIP, pod.Status.HostIP},
				IsReady:      ps.helper.PodIsReady(pod),
				EventMessage: ps.EventHandler.GetEvent(pod.Namespace, "Pod", pod.Name),
				CreateTime:   pod.CreationTimestamp.Format("2006-01-02 15:04:05.999999999 -0700 MST"),
			})
		}
	}
	return
}

// GetByLabel kubectl get pods -l app=nginx -n default
func (ps *PodService) GetByLabel(ctx context.Context, ns, label string) (*v1.PodList, error) {
	pods, err := ps.Client.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{
		LabelSelector: label,
	})
	if err != nil {
		return nil, err
	}
	return pods, nil
}
