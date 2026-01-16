package daemonset

import (
	"fmt"
	"net/http"

	"github.com/canary-debug/kube-vue-admin/pkg/global"
	"github.com/gin-gonic/gin"
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

// GetDaemonsetPods 获取 daemonset pods 请求结构体
type GetDaemonsetPods struct {
	Name      string `json:"name" binding:"required"`
	Namespace string `json:"namespace" binding:"required"`
}

func Get_Daemonset_Pods(c *gin.Context) {
	// 请求参数绑定
	var getDeploymentPods GetDaemonsetPods
	if err := c.ShouldBindJSON(&getDeploymentPods); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求体",
		})
		return
	}

	pods, err := global.Daemonsets.DaemonSets(getDeploymentPods.Namespace).Get(getDeploymentPods.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("获取Daemonset失败: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pods": pods,
	})

}
