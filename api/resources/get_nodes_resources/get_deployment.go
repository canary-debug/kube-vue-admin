package get_nodes_resources

import (
	"github.com/canary-debug/kube-vue-admin/pkg/global"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/labels"
)

var (
	deploymentName string
)

func GetDeployment(c *gin.Context) {
	// 获取命名空间
	namespace := c.Param("namespace")

	// 使用 informer 的 lister 来获取指定命名空间下的 deployments
	deploymentLister := global.Deployments

	// 通过 namespace filter 获取 deployments
	deployments, err := deploymentLister.Deployments(namespace).List(labels.Everything())
	for _, deployment := range deployments {
		// 获取 Deployment 的名字
		deploymentName = deployment.Name
	}
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"Name": deploymentName,
	})
}
