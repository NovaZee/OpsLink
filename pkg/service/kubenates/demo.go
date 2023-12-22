package kubenates

import (
	"bytes"
	"context"
	"fmt"
	"github.com/oppslink/protocol/logger"
	"io"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// todo:获取指定容器的日志,加入指定容器参数
func getLog(clientset *kubernetes.Clientset, namespace, podName string, follow bool) {
	var lines int64 = 100
	req := clientset.CoreV1().Pods(namespace).GetLogs(podName, &v1.PodLogOptions{Follow: follow, TailLines: &lines})
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		fmt.Println(err)
	}
	defer podLogs.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		fmt.Println(err)
	}
	str := buf.String()
	// 处理日志
	fmt.Println(str)
}

func getNameSpacePods(clientset *kubernetes.Clientset, namespace string) {
	get, _ := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	for _, item := range get.Items {
		logger.Infow("kubenates-system", "Namespace", item.Namespace, "Name", item.GetName())
	}
}
