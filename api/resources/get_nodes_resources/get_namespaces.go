package get_nodes_resources

import (
	"context"
	"net/http"

	"github.com/canary-debug/kube-vue-admin/api/resources"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetNameSpaces 获取所有命名空间的个数
func GetNameSpacesLen(c *gin.Context) {
	namespace, err := resources.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get namespaces: " + err.Error(),
		})
		return
	}

	// 返回 namespace 的个数
	c.JSON(http.StatusOK, gin.H{
		"namespace": len(namespace.Items),
	})

}

// 获取所有命名空间的名字
func GetNameSpaces(c *gin.Context) {
	namespace, err := resources.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get namespaces: " + err.Error(),
		})
		return
	}

	// 收集所有命名空间到切片中
	namespaces := []string{}
	for _, ns := range namespace.Items {
		namespaces = append(namespaces, ns.Name)
	}

	c.JSON(http.StatusOK, gin.H{
		"namespaces": namespaces,
	})

}
