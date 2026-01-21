package pod

import (
	"context"
	"net/http"

	"github.com/canary-debug/kube-vue-admin/api/resources"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeletePod(c *gin.Context) {
	namespace := c.Param("namespace")
	podname := c.Param("podname")

	// 获取可选的强制删除参数: /.../pod-name?force=true
	force := c.DefaultQuery("force", "false")

	var deleteOptions metav1.DeleteOptions
	if force == "true" {
		gracePeriod := int64(0)
		deleteOptions.GracePeriodSeconds = &gracePeriod
		// 强制删除通常需要配合 Orphan 策略或 Background
		policy := metav1.DeletePropagationBackground
		deleteOptions.PropagationPolicy = &policy
	}
	// 删除 Pod
	err := resources.Clientset.CoreV1().Pods(namespace).Delete(context.TODO(), podname, metav1.DeleteOptions{})
	if err != nil {
		// 如果 Pod 已经被删除了，返回 404
		if errors.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Pod 不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 204 No Content 是 DELETE 请求成功的标准返回码（代表成功且无返回主体）
	// 或者返回 200 OK 带上状态信息
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Pod " + podname + " deletion triggered",
	})

}
