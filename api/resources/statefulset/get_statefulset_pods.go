package statefulset

import (
	"fmt"
	"net/http"

	"github.com/canary-debug/kube-vue-admin/api/resources"
	"github.com/canary-debug/kube-vue-admin/api/resources/pod"
	"github.com/canary-debug/kube-vue-admin/pkg/global"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
	podList, err := global.Pods.Pods(getStatefulSetPods.Namespace).List(selector)
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
