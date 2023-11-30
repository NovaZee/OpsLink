package kubeservice

import (
	"context"
	"github.com/denovo/permission/pkg/kubenates/informer"
	"github.com/denovo/permission/protoc/kube"
	"github.com/oppslink/protocol/logger"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type DeploymentService struct {
	Di     *informer.DeploymentInformer
	Client kubernetes.Interface

	helper *helper
}

func NewDeploymentService(client kubernetes.Interface) *DeploymentService {
	return &DeploymentService{Di: &informer.DeploymentInformer{}, Client: client, helper: &helper{}}
}

func (ds *DeploymentService) List(ctx context.Context, namespace string) (res []*kube.Deployment, err error) {
	deployments, err := ds.Di.ListALl(namespace)
	if err != nil {
		return
	}

	for _, deployment := range deployments {
		res = append(res, &kube.Deployment{
			Name:       deployment.Name,
			Namespace:  deployment.Namespace,
			Replicas:   []int32{deployment.Status.Replicas, deployment.Status.AvailableReplicas, deployment.Status.UnavailableReplicas},
			Images:     ds.helper.GetImages(*deployment),
			IsComplete: ds.getDeploymentIsComplete(deployment),
			Message:    ds.getDeploymentCondition(deployment),
			CreateTime: deployment.CreationTimestamp.Format("2006-01-02 15:04:05"),
		})
	}
	return
}

func (dc *DeploymentService) Delete(ctx context.Context, ns, name string) {

	err := dc.Client.AppsV1().Deployments(ns).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		logger.Errorw("DeploymentController Delete ", err)
	}
	logger.Infow("DeploymentController Delete ", "Name", name, "Namespace", ns)
}

func (*DeploymentService) getDeploymentIsComplete(deployment *v1.Deployment) bool {
	return deployment.Status.Replicas == deployment.Status.AvailableReplicas
}

func (*DeploymentService) getDeploymentCondition(deployment *v1.Deployment) string {

	for _, item := range deployment.Status.Conditions {
		if string(item.Type) == "Available" && string(item.Status) != "True" {
			return item.Message
		}
	}
	return ""

}
