package get_nodes_resources

import (
	"context"
	"net/http"

	"github.com/canary-debug/kube-vue-admin/api/resources"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// todo: 获取所有 pod 资源
func GetContainerGroup(c *gin.Context) {
	namespace := c.Param("namespace")

	// 整合所有控制器资源
	//response := gin.H{
	//	"namespace":    namespace,
	//	"deployments":  []map[string]interface{}{},
	//	"statefulsets": []map[string]interface{}{},
	//	"daemonsets":   []map[string]interface{}{},
	//}

	/*
		获取所有 pod
		namespace 为空表示所有 pod
	*/
	pods, err := resources.Clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get pods: " + err.Error(),
		})
	}
	for _, pod := range pods.Items {
		c.JSON(http.StatusOK, gin.H{
			"pod": pod,
		})
	}

}
