package kubenates

import (
	"github.com/denovo/permission/config"
	"github.com/denovo/permission/pkg/kubenates/kubeservice"
	"github.com/oppslink/protocol/logger"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

type K8sClient struct {
	Clientset        kubernetes.Interface
	MetricsClientSet *versioned.Clientset
	RestConfig       *rest.Config

	NodeHandler      *kubeservice.NodeService
	DepHandler       *kubeservice.DeploymentService
	PodHandler       *kubeservice.PodService
	NamespaceHandler *kubeservice.NamespaceService
	EventHandler     *kubeservice.EventService
}

func NewK8sConfig(conf *config.OpsLinkConfig) (*K8sClient, error) {
	var err error
	var clientSet kubernetes.Interface
	var metricClient *versioned.Clientset
	k8sClient := &K8sClient{}

	config, err := rest.InClusterConfig()
	if err != nil {
		logger.Infow("Program running from outside of the cluster")
		set, mc, err2 := NewClientSet(conf)
		if err2 != nil {
			return nil, err2
		}
		clientSet = set
		metricClient = mc
		err = nil
	} else {
		kubeConfig :=
			clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
		set, err2 := kubernetes.NewForConfig(config)
		if err2 != nil {
			return nil, err2
		}
		clientSet = set

		mc, err3 := versioned.NewForConfig(config)
		if err3 != nil {
			return nil, err3
		}
		metricClient = mc
	}
	if err != nil {
		logger.Infow("Program running from outside of the cluster")

	} else {
		logger.Infow("Program running inside the cluster, picking the in-cluster configuration")
	}

	k8sClient.Clientset = clientSet
	k8sClient.MetricsClientSet = metricClient
	k8sClient.RestConfig = config

	//init resource
	k8sClient.initHandlers()

	k8sClient.InitInformer()

	return k8sClient, err
}

// initHandlers 用于初始化 DepHandler 和 PodHandler
func (k *K8sClient) initHandlers() {
	k.EventHandler = kubeservice.NewEventService(k.Clientset)

	k.DepHandler = kubeservice.NewDeploymentService(k.Clientset, k.EventHandler)
	k.PodHandler = kubeservice.NewPodService(k.Clientset, k.EventHandler)
	k.NamespaceHandler = kubeservice.NewNamespaceService(k.Clientset)

	k.NodeHandler = kubeservice.NewNodeService(k.Clientset, k.MetricsClientSet, k.PodHandler.Pi)

}

// NewClientSet Kubernetes客户端的接口实例D
func NewClientSet(conf *config.OpsLinkConfig) (kubernetes.Interface, *versioned.Clientset, error) {
	var err error
	kubeconfig := conf.Kubernetes.Kubeconfig
	configOverrides := &clientcmd.ConfigOverrides{}
	var kubecfg *rest.Config

	kubecfg, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		configOverrides).ClientConfig()
	if err != nil {
		return nil, nil, err
	}
	k8sClient, err := kubernetes.NewForConfig(kubecfg)
	if err != nil {
		return nil, nil, err
	}

	mc, err3 := versioned.NewForConfig(kubecfg)
	if err3 != nil {
		return k8sClient, nil, err3
	}

	return k8sClient, mc, nil
}

// InitInformer informer初始化
func (k *K8sClient) InitInformer() informers.SharedInformerFactory {
	sif := informers.NewSharedInformerFactory(k.Clientset, 0)

	deploymentInformer := sif.Apps().V1().Deployments()
	deploymentInformer.Informer().AddEventHandler(k.DepHandler.Di)

	node := sif.Core().V1().Nodes()
	node.Informer().AddEventHandler(k.NodeHandler.Ni)

	pods := sif.Core().V1().Pods()
	pods.Informer().AddEventHandler(k.PodHandler.Pi)

	event := sif.Core().V1().Events()
	event.Informer().AddEventHandler(k.EventHandler.Ei)

	ns := sif.Core().V1().Namespaces()
	ns.Informer().AddEventHandler(k.NamespaceHandler.Nsi)

	sif.Start(wait.NeverStop)

	return sif

}
