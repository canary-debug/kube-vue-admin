package deployment

import (
	"fmt"
	"net/http"

	"github.com/canary-debug/kube-vue-admin/api/resources"
	"github.com/canary-debug/kube-vue-admin/pkg/global"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PodInfo Pod信息结构体
type PodInfo struct {
	Name         string            `json:"name"`
	Status       string            `json:"status"`
	RestartCount int32             `json:"restart_count"`
	Ports        []int32           `json:"ports"`
	NodeName     string            `json:"node_name"`
	PodIP        string            `json:"pod_ip"`
	CreatedAt    metav1.Time       `json:"created_at"`
	Labels       map[string]string `json:"labels"`
}

// GetDeploymentPods 获取 deployment pods 请求结构体
type GetDeploymentPods struct {
	Name      string `json:"name" binding:"required"`
	Namespace string `json:"namespace" binding:"required"`
}

func Get_Deployment_Pods(c *gin.Context) {
	// 请求参数绑定
	var getDeploymentPods GetDeploymentPods
	if err := c.ShouldBindJSON(&getDeploymentPods); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求体",
		})
		return
	}

	// 验证ClientSet是否已初始化
	if resources.Clientset == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Kubernetes客户端未初始化",
		})
		return
	}

	// 验证全局Informer是否已初始化
	if global.Deployments == nil || global.Pods == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "全局Informer未初始化",
		})
		return
	}

	// 获取Deployment信息以获取标签选择器
	deployment, err := global.Deployments.Deployments(getDeploymentPods.Namespace).Get(getDeploymentPods.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("获取Deployment失败: %v", err),
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("获取Deployment失败: %v", err),
		})
		return
	}

	// 使用Deployment的标签选择器筛选Pod
	selector, err := metav1.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("转换标签选择器失败: %v", err),
		})
		return
	}

	// 使用全局Pod Lister获取Pod列表
	pods, err := global.Pods.Pods(getDeploymentPods.Namespace).List(selector)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("获取Pod列表失败: %v", err),
		})
		return
	}

	// 处理Pod信息并构建返回结果
	var podInfos []PodInfo
	for _, pod := range pods {
		podInfo := convertPodToPodInfo(pod)
		podInfos = append(podInfos, *podInfo)
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"data":    podInfos,
		"total":   len(podInfos),
		"message": "获取成功",
	})

}

// convertPodToPodInfo 将Pod对象转换为PodInfo结构体
func convertPodToPodInfo(pod *v1.Pod) *PodInfo {
	// 计算总重启次数
	totalRestartCount := int32(0)
	for _, status := range pod.Status.InitContainerStatuses {
		totalRestartCount += status.RestartCount
	}
	for _, status := range pod.Status.ContainerStatuses {
		totalRestartCount += status.RestartCount
	}

	// 获取容器端口
	var ports []int32
	for _, container := range pod.Spec.Containers {
		for _, port := range container.Ports {
			ports = append(ports, port.ContainerPort)
		}
	}

	// 获取Pod状态
	status := getPodStatus(pod)

	return &PodInfo{
		Name:         pod.Name,
		Status:       status,
		RestartCount: totalRestartCount,
		Ports:        ports,
		NodeName:     pod.Spec.NodeName,
		PodIP:        pod.Status.PodIP,
		CreatedAt:    pod.CreationTimestamp,
		Labels:       pod.Labels,
	}
}

// getPodStatus 获取Pod的简化状态
func getPodStatus(pod *v1.Pod) string {
	switch pod.Status.Phase {
	case v1.PodPending:
		return "Pending"
	case v1.PodRunning:
		// 检查是否有容器处于等待或终止状态
		for _, status := range pod.Status.ContainerStatuses {
			if status.State.Waiting != nil {
				return "Waiting"
			}
			if status.State.Terminated != nil {
				return "Terminated"
			}
		}
		return "Running"
	case v1.PodSucceeded:
		return "Succeeded"
	case v1.PodFailed:
		return "Failed"
	case v1.PodUnknown:
		return "Unknown"
	default:
		return string(pod.Status.Phase)
	}
}
