package clientset

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var kubeconfig *string
	//if home := homedir.HomeDir(); home != "" {
	//	kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kubeclient", "config"), "(optional) absolute path to the kubeconfig file")
	//} else {
	//	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	//}
	//flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	// 获取mysql的pod
	namespace := "whalebase"
	pod := "mysql-0"
	mysqlPod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), pod, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting pod %s in namespace %s: %v\n",
			pod, namespace, statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
		bytes, _ := json.Marshal(mysqlPod)
		fmt.Println("pod信息：", string(bytes))
	}

	fmt.Println("--------------------------------")
	// 获取mysql的StatefulSets
	sts, _ := clientset.AppsV1().StatefulSets(namespace).Get(context.TODO(), "mysql", metav1.GetOptions{})
	bytes, _ := json.Marshal(sts)
	fmt.Println("sts信息：", string(bytes))
	fmt.Println("--------------------------------")
	//创建pod
	var nginxPod *v1.Pod = &v1.Pod{
		TypeMeta:   metav1.TypeMeta{Kind: "Pod", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{Name: "nginx-pod"},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				v1.Container{
					Name:  "nginx",
					Image: "nginx:1.8",
				},
			},
		},
	}
	_, err = clientset.CoreV1().Pods("default").Create(context.TODO(), nginxPod, metav1.CreateOptions{})
	if err != nil {
		fmt.Println(err)
	}
}
