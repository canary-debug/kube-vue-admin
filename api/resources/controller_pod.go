package resources

import (
	"net/http"

	"github.com/canary-debug/kube-vue-admin/pkg/global"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/labels"
)

func GetPod(c *gin.Context) {
	// 获取命名空间, 资源名称 返回 pod 信息
	namespace := c.Param("namespace")
	resourcename := c.Param("resourcename")
	if namespace == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "namespace parameter is required",
		})
		return
	} else if resourcename == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "resourcename parameter is required",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"resourcename": resourcename,
		"namespace":    namespace,
	})

}

// GetPodCount 返回集群所有 Pod 数量
func GetPodCount(c *gin.Context) {
	// 使用 informer 的 lister 来获取 Pod 总数量
	podCount := global.Pods
	podLen, err := podCount.List(labels.Everything())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get pod count",
		})
	}

	// 返回 Pod 数量
	c.JSON(http.StatusOK, gin.H{
		"pod_count": len(podLen),
	})

}
