package infomer

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

// 获取系统家目录
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}

	// for windows
	return os.Getenv("USERPROFILE")
}

func main() {
	var kubeConfig *string

	if h := homeDir(); h != "" {
		kubeConfig = flag.String("kubeconfig", filepath.Join(h, ".kubeclient", "config"), "use kubeconfig access to kubeapiserver")
	} else {
		kubeConfig = flag.String("kubeconfig", "", "use kubeconfig access to kubeapiserver")
	}

	// 获取 kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if err != nil {
		panic(err.Error())
	}
	// 初始化 clientSet
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// 初始化 default 命名空间下的 informer 工厂, 这个 informer 工厂包含 k8s 所有内置资源的 informer
	// 同时设置 5s 的同步周期，同步是指将 indexer 的数据同步到 deltafifo，防止因为特殊原因处理失败的数据能够得到重新处理
	informerFactory := informers.NewSharedInformerFactoryWithOptions(clientSet, 5*time.Second, informers.WithNamespace("default"))
	// 获取 pod informer
	podInformer := informerFactory.Core().V1().Pods().Informer()
	// 向 pod informer 注册处理函数
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    add,
		UpdateFunc: update,
		DeleteFunc: delete,
	})
	//获取hpa infomer
	//autoscalers := informerFactory.Autoscaling().V1().HorizontalPodAutoscalers()
	//autoscalers.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
	//	AddFunc:    add,
	//	UpdateFunc: update,
	//	DeleteFunc: delete,
	//})

	stopChan := make(chan struct{})
	defer close(stopChan)
	// 启动 pod informer
	podInformer.Run(stopChan)
	// 等待数据同步到 cache 中
	isCache := cache.WaitForCacheSync(stopChan, podInformer.HasSynced)
	if !isCache {
		fmt.Println("pod has not cached")
		return
	}
}

// 资源新增回调函数
func add(obj interface{}) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		panic("invalid obj")
	}
	fmt.Println("add a pod:", pod.Name)
}

// 资源更新回调函数
func update(oldObj, newObj interface{}) {
	oldPod, ok := oldObj.(*corev1.Pod)
	if !ok {
		panic("invalid oldObj")
	}
	newPod, ok := newObj.(*corev1.Pod)
	if !ok {
		panic("invalid newObj")
	}
	fmt.Println("update a pod:", oldPod.Name, newPod.Name)
}

// 资源删除回调函数
func delete(obj interface{}) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		panic("invalid obj")
	}
	fmt.Println("delete a pod:", pod.Name)
}

//todo:https://github.com/rancher/lasso 控制器框架
//todo：https://kubernetes.io/docs/concepts/overview/components/#addons
