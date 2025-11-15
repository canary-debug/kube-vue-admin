package api

import (
	"context"
	"flag"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"net/http"
	"path/filepath"
	"time"
)

var Clientset *kubernetes.Clientset

// 初始化 Kubernetes 客户端
// 创建 Clientset
func init() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	Clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
}

// GetControllersInNamespace 获取指定命名空间下的所有控制器资源
func GetControllersInNamespace(c *gin.Context) {
	namespace := c.Param("namespace")
	if namespace == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "namespace parameter is required",
		})
		return
	}

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 获取 Deployments
	deployments, err := Clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get deployments: " + err.Error(),
		})
		return
	}

	// 获取 StatefulSets
	statefulSets, err := Clientset.AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get statefulsets: " + err.Error(),
		})
		return
	}

	// 获取 DaemonSets
	daemonSets, err := Clientset.AppsV1().DaemonSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get daemonsets: " + err.Error(),
		})
		return
	}

	// 整合所有控制器资源
	response := gin.H{
		"namespace":    namespace,
		"deployments":  []map[string]interface{}{},
		"statefulsets": []map[string]interface{}{},
		"daemonsets":   []map[string]interface{}{},
	}

	// 处理 Deployments
	for _, deployment := range deployments.Items {
		response["deployments"] = append(response["deployments"].([]map[string]interface{}), map[string]interface{}{
			"name":       deployment.Name,
			"replicas":   deployment.Spec.Replicas,
			"ready":      deployment.Status.ReadyReplicas,
			"updated":    deployment.Status.UpdatedReplicas,
			"available":  deployment.Status.AvailableReplicas,
			"created_at": deployment.CreationTimestamp,
		})
	}

	// 处理 StatefulSets
	for _, statefulSet := range statefulSets.Items {
		response["statefulsets"] = append(response["statefulsets"].([]map[string]interface{}), map[string]interface{}{
			"name":       statefulSet.Name,
			"replicas":   statefulSet.Spec.Replicas,
			"ready":      statefulSet.Status.ReadyReplicas,
			"updated":    statefulSet.Status.UpdatedReplicas,
			"available":  statefulSet.Status.CurrentReplicas,
			"created_at": statefulSet.CreationTimestamp,
		})
	}

	// 处理 DaemonSets
	for _, daemonSet := range daemonSets.Items {
		response["daemonsets"] = append(response["daemonsets"].([]map[string]interface{}), map[string]interface{}{
			"name":       daemonSet.Name,
			"desired":    daemonSet.Status.DesiredNumberScheduled,
			"current":    daemonSet.Status.CurrentNumberScheduled,
			"ready":      daemonSet.Status.NumberReady,
			"updated":    daemonSet.Status.UpdatedNumberScheduled,
			"available":  daemonSet.Status.NumberAvailable,
			"created_at": daemonSet.CreationTimestamp,
		})
	}

	c.JSON(http.StatusOK, response)
}
