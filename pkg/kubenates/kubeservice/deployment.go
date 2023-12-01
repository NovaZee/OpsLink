package kubeservice

import (
	"context"
	"github.com/denovo/permission/pkg/kubenates/informer"
	"github.com/denovo/permission/protoc/kube"
	"github.com/oppslink/protocol/logger"
	"gopkg.in/yaml.v3"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type DeploymentService struct {
	Di     *informer.DeploymentInformer
	Client kubernetes.Interface

	helper *helper
}

type StructMetadata struct {
	Annotations       map[string]string `json:"annotations"`
	CreationTimestamp string            `json:"creationTimestamp"`
	Generation        int64             `json:"generation"`
	Labels            map[string]string `json:"labels"`
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	ResourceVersion   string            `json:"resourceVersion"`
	UID               string            `json:"uid"`
}

func NewDeploymentService(client kubernetes.Interface) *DeploymentService {
	return &DeploymentService{Di: &informer.DeploymentInformer{}, Client: client, helper: &helper{}}
}

func (ds *DeploymentService) List(namespace string) (res []*kube.Deployment, err error) {
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

func (ds *DeploymentService) DownToYaml(ns, name string) ([]byte, error) {
	deployments, err := ds.Di.ListALl(ns)
	if err != nil {
		return nil, err
	}
	for _, deployment := range deployments {
		if name == deployment.Name {

			// Create a struct to hold partial deployment information including apiVersion, kind, metadata, spec, and status
			metadate := &StructMetadata{
				Annotations:       deployment.ObjectMeta.Annotations,
				CreationTimestamp: deployment.CreationTimestamp.String(),
				Generation:        deployment.Generation,
				Labels:            deployment.Labels,
				Name:              deployment.Name,
				Namespace:         deployment.Namespace,
				ResourceVersion:   deployment.ResourceVersion,
				UID:               string(deployment.UID),
			}
			partialDeployment := struct {
				APIVersion string              `json:"apiVersion"`
				Kind       string              `json:"kind"`
				Metadata   *StructMetadata     `json:"metadata"`
				Spec       v1.DeploymentSpec   `json:"spec,omitempty"`
				Status     v1.DeploymentStatus `json:"status,omitempty"`
			}{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Metadata:   metadate,
				Spec:       deployment.Spec,
				Status:     deployment.Status,
			}

			deploymentByte, err := yaml.Marshal(partialDeployment)
			if err != nil {
				return nil, err
			}
			return deploymentByte, nil
		}
	}
	return nil, nil
}

func (ds *DeploymentService) Apply(ns string, deployment *v1.Deployment) error {

	_, err := ds.Client.AppsV1().Deployments(ns).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (ds *DeploymentService) Delete(ctx context.Context, ns, name string) error {

	err := ds.Client.AppsV1().Deployments(ns).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return err
		logger.Errorw("DeploymentController Delete ", err)
	}
	logger.Infow("DeploymentController Delete ", "Name", name, "Namespace", ns)
	return nil
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
