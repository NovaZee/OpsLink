package kubeservice

import (
	"context"
	"github.com/denovo/permission/pkg/service/kubenates/informer"
	"github.com/denovo/permission/protoc/kube"
	"github.com/oppslink/protocol/logger"
	v3yaml "gopkg.in/yaml.v3"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

type DeploymentService struct {
	Di     *informer.DeploymentInformer
	Client kubernetes.Interface

	EventHandler *EventService
	helper       *helper
}

func NewDeploymentService(client kubernetes.Interface, eh *EventService) *DeploymentService {
	return &DeploymentService{Di: &informer.DeploymentInformer{}, Client: client, helper: &helper{}, EventHandler: eh}
}

func (ds *DeploymentService) GetDeployment(ns, name string) (*v1.Deployment, error) {
	deployments, err := ds.Di.ListALl(ns)
	if err != nil {
		return nil, err
	}
	for _, deployment := range deployments {
		if name == deployment.Name {
			return deployment, nil
		}
	}
	return nil, nil
}

func (ds *DeploymentService) List(namespace string) (res []*kube.Deployment, err error) {
	deployments, err := ds.Di.ListALl(namespace)
	if err != nil {
		return
	}

	for _, deployment := range deployments {
		res = append(res, &kube.Deployment{
			Name:         deployment.Name,
			Namespace:    deployment.Namespace,
			Replicas:     []int32{deployment.Status.Replicas, deployment.Status.AvailableReplicas, deployment.Status.UnavailableReplicas},
			Images:       ds.helper.GetImages(*deployment),
			IsComplete:   ds.getDeploymentIsComplete(deployment),
			Message:      ds.getDeploymentCondition(deployment),
			EventMessage: ds.EventHandler.GetEvent(namespace, "Deployment", deployment.Name),
			CreateTime:   deployment.CreationTimestamp.Format("2006-01-02 15:04:05"),
			Labels:       ds.helper.LabelsFilter(deployment.Spec.Selector.MatchLabels),
		})
	}
	return
}

func (ds *DeploymentService) Update(ctx context.Context, ns, name string, updateData *v1.Deployment) ([]byte, error) {
	updatedDeployment, err := ds.Client.AppsV1().Deployments(ns).Update(ctx, updateData, metav1.UpdateOptions{})
	if err != nil {
		logger.Warnw("DeploymentController Update ", err, "Name", name, "Namespace", ns)
		return nil, err
	}

	deploymentData, converterErr := runtime.DefaultUnstructuredConverter.ToUnstructured(updatedDeployment)
	if converterErr != nil {
		return nil, converterErr
	}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(deploymentData, &v1.Deployment{})
	if err != nil {
		return nil, err
	}
	deploymentData["apiVersion"] = "apps/v1"
	deploymentData["kind"] = "Deployment"
	deploymentByte, err := v3yaml.Marshal(deploymentData)
	if err != nil {
		return nil, err
	}
	return deploymentByte, nil
}

func (ds *DeploymentService) Patch(ctx context.Context, ns, name string, patchData []byte) (res *v1.Deployment, err error) {
	//四种patch，默认选择json_patch,后续增加其他
	updatedDeployment, err := ds.Client.AppsV1().Deployments(ns).Patch(ctx, name, types.JSONPatchType, patchData, metav1.PatchOptions{})
	if err != nil {
		logger.Warnw("DeploymentController Patch ", err, "Name", name, "Namespace", ns)
	}
	return updatedDeployment, err
}

func (ds *DeploymentService) DownToYaml(ns, name string) ([]byte, error) {
	deployments, err := ds.Di.ListALl(ns)
	if err != nil {
		return nil, err
	}
	for _, deployment := range deployments {
		if name == deployment.Name {
			//将 Deployment 对象转换为 Unstructured 对象
			deploymentData, converterErr := runtime.DefaultUnstructuredConverter.ToUnstructured(deployment)
			if converterErr != nil {
				return nil, converterErr
			}
			err = runtime.DefaultUnstructuredConverter.FromUnstructured(deploymentData, &v1.Deployment{})
			if err != nil {
				return nil, err
			}
			deploymentData["apiVersion"] = "apps/v1"
			deploymentData["kind"] = "Deployment"
			deploymentByte, err := v3yaml.Marshal(deploymentData)
			if err != nil {
				return nil, err
			}
			return deploymentByte, nil
		}
	}
	return nil, nil
}

func (ds *DeploymentService) ApplyByYaml(ctx context.Context, ns string, in []byte, isUpdate bool) error {
	// create unstructured object
	decode := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	deployment := &v1.Deployment{}
	_, _, err := decode.Decode(in, nil, deployment)
	if err != nil {
		return err
	}
	if isUpdate {
		_, err = ds.Client.AppsV1().Deployments(ns).Update(ctx, deployment, metav1.UpdateOptions{})
	} else {
		_, err = ds.Client.AppsV1().Deployments(ns).Create(ctx, deployment, metav1.CreateOptions{})
	}
	if err != nil {
		return err
	}
	return nil
}

func (ds *DeploymentService) Rollout(ctx context.Context, ns, name string) error {
	deployment, err := ds.GetDeployment(ns, name)
	if err != nil {
		return err
	}
	// 更新 Deployment 的标签
	deployment.Spec.Template.ObjectMeta.Labels["restartedAt"] = metav1.Now().Format("2006-01-02_15-04-05")

	_, updateErr := ds.Client.AppsV1().Deployments(ns).Update(ctx, deployment, metav1.UpdateOptions{})
	if updateErr != nil {
		return updateErr
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
