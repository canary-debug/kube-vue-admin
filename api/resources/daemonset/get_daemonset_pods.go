package daemonset

import (
	"fmt"
	"net/http"

	"github.com/canary-debug/kube-vue-admin/api/resources"
	"github.com/canary-debug/kube-vue-admin/api/resources/pod"
	"github.com/canary-debug/kube-vue-admin/pkg/global"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetDaemonSetPods 获取 DaemonSet pods 请求结构体
type GetDaemonSetPods struct {
	Name      string `json:"name" binding:"required"`
	Namespace string `json:"namespace" binding:"required"`
}

// GetDaemonsetPods 获取指定 DaemonSet 下的所有 Pod
func GetDaemonsetPods(c *gin.Context) {
	// 1. 请求参数绑定
	var getDaemonSetPods GetDaemonSetPods
	if err := c.ShouldBindJSON(&getDaemonSetPods); err != nil {
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

	if global.Daemonsets == nil || global.Pods == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "全局DaemonSet或Pod Informer未初始化",
		})
		return
	}

	// 3. 获取 DaemonSet 信息
	ds, err := global.Daemonsets.DaemonSets(getDaemonSetPods.Namespace).Get(getDaemonSetPods.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("获取DaemonSet失败: %v", err),
		})
		return
	}

	// 4. 使用 DaemonSet 的标签选择器转换成 Selector 对象
	selector, err := metav1.LabelSelectorAsSelector(ds.Spec.Selector)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("转换标签选择器失败: %v", err),
		})
		return
	}

	// 5. 使用 Pod Lister 筛选符合该 Selector 的 Pod
	podList, err := global.Pods.Pods(getDaemonSetPods.Namespace).List(selector)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("获取Pod列表失败: %v", err),
		})
		return
	}

	// 6. 处理 Pod 信息并构建返回结果
	var podInfos []pod.PodInfo
	for _, p := range podList {
		podInfo := pod.ConvertPodToPodInfo(p)
		podInfos = append(podInfos, *podInfo)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    podInfos,
		"total":   len(podInfos),
		"message": "获取成功",
	})
}
