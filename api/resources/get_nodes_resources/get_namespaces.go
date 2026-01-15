package get_nodes_resources

import (
	"context"
	"log"
	"net/http"

	"github.com/canary-debug/kube-vue-admin/api/resources"
	"github.com/canary-debug/kube-vue-admin/pkg/global"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// GetNameSpacesLen GetNameSpaces 获取所有命名空间的个数
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
	NameSpace := global.Namespaces
	if NameSpace == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Informer 尚未初始化完成，请稍后重试"})
		return
	}
	// 从本地 Lister 中获取所有 Namespace，不直接请求 API Server
	//NameSpace := global.Namespaces
	namespaces, err := NameSpace.List(labels.Everything())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var names []string
	for _, ns := range namespaces {
		names = append(names, ns.Name)
	}
	log.Println("获取所有命名空间名字: ", names)

	c.JSON(http.StatusOK, gin.H{
		//"count": len(names),
		"namespaces": names,
	})

}
