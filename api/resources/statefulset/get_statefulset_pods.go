package statefulset

import (
	"fmt"
	"net/http"

	"github.com/canary-debug/kube-vue-admin/api/resources"
	"github.com/canary-debug/kube-vue-admin/pkg/global"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PodInfo Pod信息结构体 (保持不变)
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

// GetStatefulSetPods 获取 StatefulSet pods 请求结构体
type GetStatefulSetPods struct {
	Name      string `json:"name" binding:"required"`
	Namespace string `json:"namespace" binding:"required"`
}

// GetStatefulsetPods 获取指定 StatefulSet 下的所有 Pod
func GetStatefulsetPods(c *gin.Context) {
	// 1. 请求参数绑定
	var getStatefulSetPods GetStatefulSetPods
	if err := c.ShouldBindJSON(&getStatefulSetPods); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求体",
		})
		return
	}

	// 2. 验证客户端和 Informer 是否初始化
	if resources.Clientset == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Kubernetes客户端未初始化",
		})
		return
	}

	// 注意：确保 global 包中定义并初始化了 Statefulsets 的 Informer/Lister
	if global.Statefulset == nil || global.Pods == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "全局StatefulSet或Pod Informer未初始化",
		})
		return
	}

	// 3. 获取 StatefulSet 信息 (用于获取 LabelSelector)
	sts, err := global.Statefulset.StatefulSets(getStatefulSetPods.Namespace).Get(getStatefulSetPods.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("获取StatefulSet失败: %v", err),
		})
		return
	}

	// 4. 使用 StatefulSet 的标签选择器转换成 Selector 对象
	selector, err := metav1.LabelSelectorAsSelector(sts.Spec.Selector)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("转换标签选择器失败: %v", err),
		})
		return
	}

	// 5. 使用 Pod Lister 筛选符合该 Selector 的 Pod
	pods, err := global.Pods.Pods(getStatefulSetPods.Namespace).List(selector)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("获取Pod列表失败: %v", err),
		})
		return
	}

	// 6. 处理 Pod 信息并构建返回结果
	var podInfos []PodInfo
	for _, pod := range pods {
		podInfo := convertPodToPodInfo(pod)
		podInfos = append(podInfos, *podInfo)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    podInfos,
		"total":   len(podInfos),
		"message": "获取成功",
	})
}

// convertPodToPodInfo (保持原有逻辑不变)
func convertPodToPodInfo(pod *v1.Pod) *PodInfo {
	totalRestartCount := int32(0)
	for _, status := range pod.Status.InitContainerStatuses {
		totalRestartCount += status.RestartCount
	}
	for _, status := range pod.Status.ContainerStatuses {
		totalRestartCount += status.RestartCount
	}

	var ports []int32
	for _, container := range pod.Spec.Containers {
		for _, port := range container.Ports {
			ports = append(ports, port.ContainerPort)
		}
	}

	return &PodInfo{
		Name:         pod.Name,
		Status:       getPodStatus(pod),
		RestartCount: totalRestartCount,
		Ports:        ports,
		NodeName:     pod.Spec.NodeName,
		PodIP:        pod.Status.PodIP,
		CreatedAt:    pod.CreationTimestamp,
		Labels:       pod.Labels,
	}
}

// getPodStatus (保持原有逻辑不变)
func getPodStatus(pod *v1.Pod) string {
	switch pod.Status.Phase {
	case v1.PodPending:
		return "Pending"
	case v1.PodRunning:
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
	default:
		return string(pod.Status.Phase)
	}
}
