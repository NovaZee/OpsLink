package kubeservice

import (
	"context"
	"github.com/denovo/permission/pkg/kubenates/informer"
	"github.com/denovo/permission/protoc/kube"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

type NodeService struct {
	Ni     *informer.NodeInformer
	Pi     *informer.PodInformer
	Metric *versioned.Clientset
	Client kubernetes.Interface
	helper *helper
}

func NewNodeService(client kubernetes.Interface, mc *versioned.Clientset, pi *informer.PodInformer) *NodeService {
	return &NodeService{Ni: &informer.NodeInformer{}, Pi: pi, Client: client, Metric: mc, helper: &helper{}}
}
func (ns *NodeService) Get(nodeName string) *corev1.Node {
	return ns.Ni.Get(nodeName)
}

func (ns *NodeService) Update(ctx context.Context, node *corev1.Node) (*corev1.Node, error) {
	node, err := ns.Client.CoreV1().Nodes().Update(ctx, node, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (ns *NodeService) List(ctx context.Context) (res []*kube.Node) {
	all := ns.Ni.ListAll()
	for _, node := range all {
		usage := GetUsage(ns.Metric, node, ctx)
		res = append(res, &kube.Node{
			Name:     node.Name,
			Ip:       node.Status.Addresses[0].Address,
			HostName: node.Status.Addresses[1].Address,
			Labels:   ns.helper.LabelsFilter(node.Labels),
			Taints:   ns.helper.TaintsFilter(node.Spec.Taints),
			Capacity: &kube.NodeCapacity{
				Cpu:    node.Status.Capacity.Cpu().Value(),
				Memory: node.Status.Capacity.Memory().Value(),
				Pods:   node.Status.Capacity.Pods().Value(),
			},
			Usage: &kube.NodeUsage{
				Pods:   int32(ns.GetPodsNum(node.Name)),
				Cpu:    usage[0],
				Memory: usage[1],
			},
			CreateTime: node.CreationTimestamp.Format("2006-01-02 15:04:05"),
		})
	}
	return
}

func GetUsage(c *versioned.Clientset, node *v1.Node, ctx context.Context) []float64 {
	nodeMetric, _ := c.MetricsV1beta1().
		NodeMetricses().Get(ctx, node.Name, metav1.GetOptions{})
	cpu := float64(nodeMetric.Usage.Cpu().MilliValue()) / float64(node.Status.Capacity.Cpu().MilliValue())
	memory := float64(nodeMetric.Usage.Memory().MilliValue()) / float64(node.Status.Capacity.Memory().MilliValue())
	return []float64{cpu, memory}
}

// GetPodsNum 根据节点名称 获取pods数量
func (ns *NodeService) GetPodsNum(node string) (num int) {
	pods, _ := ns.Pi.ListALl()
	for _, pod := range pods {
		if pod.Spec.NodeName == node {
			num++
		}
	}
	return
}
