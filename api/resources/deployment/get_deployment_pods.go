package deployment

import (
	"fmt"
	"net/http"

	"github.com/canary-debug/kube-vue-admin/api/resources"
	"github.com/canary-debug/kube-vue-admin/api/resources/pod"
	"github.com/canary-debug/kube-vue-admin/pkg/global"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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

	// 使用Deployment的标签选择器筛选Pod
	selector, err := metav1.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("转换标签选择器失败: %v", err),
		})
		return
	}

	// 使用全局Pod Lister获取Pod列表
	podList, err := global.Pods.Pods(getDeploymentPods.Namespace).List(selector)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("获取Pod列表失败: %v", err),
		})
		return
	}

	// 处理Pod信息并构建返回结果
	var podInfos []pod.PodInfo
	for _, p := range podList {
		podInfo := pod.ConvertPodToPodInfo(p)
		podInfos = append(podInfos, *podInfo)
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"data":    podInfos,
		"total":   len(podInfos),
		"message": "获取成功",
	})

}
