package resources

import (
	"context"
	"flag"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var Clientset *kubernetes.Clientset

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
	// 获取命名空间
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
		// 从 pod 中获取所有 containerd 镜像
		var deploymentImages []string
		for _, container := range deployment.Spec.Template.Spec.Containers {
			deploymentImages = append(deploymentImages, container.Image)
		}

		// 安全地获取端口, 没有端口将返回 0
		var port int32
		if len(deployment.Spec.Template.Spec.Containers) > 0 &&
			len(deployment.Spec.Template.Spec.Containers[0].Ports) > 0 {
			port = deployment.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort
		}

		response["deployments"] = append(response["deployments"].([]map[string]interface{}), map[string]interface{}{
			"name":       deployment.Name,
			"replicas":   deployment.Spec.Replicas,
			"images":     deploymentImages,
			"ready":      deployment.Status.ReadyReplicas,
			"updated":    deployment.Status.UpdatedReplicas,
			"available":  deployment.Status.AvailableReplicas,
			"created_at": deployment.CreationTimestamp,
			"update_at":  deployment.CreationTimestamp.Time,
			"port":       port,
		})
	}

	// 处理 StatefulSets
	for _, statefulSet := range statefulSets.Items {
		// 从 pod 中获取所有 containerd 镜像
		var statefulSetImages []string
		for _, container := range statefulSet.Spec.Template.Spec.Containers {
			statefulSetImages = append(statefulSetImages, container.Image)
		}

		// 安全地获取端口, 没有端口将返回 0
		var port int32
		if len(statefulSet.Spec.Template.Spec.Containers) > 0 &&
			len(statefulSet.Spec.Template.Spec.Containers[0].Ports) > 0 {
			port = statefulSet.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort
		}

		response["statefulsets"] = append(response["statefulsets"].([]map[string]interface{}), map[string]interface{}{
			"name":       statefulSet.Name,
			"replicas":   statefulSet.Spec.Replicas,
			"images":     statefulSetImages,
			"ready":      statefulSet.Status.ReadyReplicas,
			"updated":    statefulSet.Status.UpdatedReplicas,
			"available":  statefulSet.Status.CurrentReplicas,
			"created_at": statefulSet.CreationTimestamp,
			"update_at":  statefulSet.CreationTimestamp.Time,
			"port":       port,
		})
	}

	// 处理 DaemonSets
	for _, daemonSet := range daemonSets.Items {
		// 从 pod 中获取所有 containerd 镜像
		var daemonSetImages []string
		for _, container := range daemonSet.Spec.Template.Spec.Containers {
			daemonSetImages = append(daemonSetImages, container.Image)
		}

		// 安全地获取端口, 没有端口将返回 0
		var port int32
		if len(daemonSet.Spec.Template.Spec.Containers) > 0 &&
			len(daemonSet.Spec.Template.Spec.Containers[0].Ports) > 0 {
			port = daemonSet.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort
		}

		response["daemonsets"] = append(response["daemonsets"].([]map[string]interface{}), map[string]interface{}{
			"name":       daemonSet.Name,
			"images":     daemonSetImages,
			"desired":    daemonSet.Status.DesiredNumberScheduled,
			"current":    daemonSet.Status.CurrentNumberScheduled,
			"ready":      daemonSet.Status.NumberReady,
			"updated":    daemonSet.Status.UpdatedNumberScheduled,
			"available":  daemonSet.Status.NumberAvailable,
			"created_at": daemonSet.CreationTimestamp,
			"update_at":  daemonSet.CreationTimestamp.Time,
			"port":       port,
		})
	}

	c.JSON(http.StatusOK, response)
}
