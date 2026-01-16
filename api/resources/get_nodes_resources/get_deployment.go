package get_nodes_resources

import (
	"fmt"
	"log"

	"github.com/canary-debug/kube-vue-admin/pkg/global"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/labels"
)

func GetDeployment(c *gin.Context) {
	// 创建一个切片，用于存储 deployment 的名字和状态信息
	deploymentStatus := []interface{}{}

	// 获取命名空间
	namespace := c.Param("namespace")
	log.Println("namespace: ", namespace)

	// 使用 informer 的 lister 来获取指定命名空间下的 deployments
	deploymentLister := global.Deployments

	// 通过 namespace filter 获取 deployments
	deployments, err := deploymentLister.Deployments(namespace).List(labels.Everything())
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	// 遍历 deployments 返回
	for _, deployment := range deployments {
		// 获取 deployment 的名字
		//deploymentName = append(deploymentName, deployment.Name)
		// 获取 deployment 的状态信息
		statusInfo := map[string]interface{}{
			"name":     deployment.Name,
			"replicas": deployment.Status.Replicas,
		}

		// 获取运行状态
		var runningStatus string
		for _, condition := range deployment.Status.Conditions {
			if condition.Type == "Available" {
				if condition.Status == "True" {
					// 使用格式化字符串构建 (实际/期望) 格式
					runningStatus = fmt.Sprintf("运行中 (%d/%d)",
						deployment.Status.ReadyReplicas,
						*deployment.Spec.Replicas)
				} else {
					runningStatus = fmt.Sprintf("未就绪 (%d/%d)",
						deployment.Status.ReadyReplicas,
						*deployment.Spec.Replicas)
				}
				break
			}
		}
		statusInfo["status"] = runningStatus

		// 从 ManagedFields 获取最后更新时间
		if len(deployment.ManagedFields) > 0 {
			lastField := deployment.ManagedFields[len(deployment.ManagedFields)-1]
			if lastField.Time != nil {
				// 格式化时间，去除时区信息
				formattedTime := lastField.Time.Format("2006-01-02 15:04:05")
				statusInfo["update_time"] = formattedTime
			}
		}

		deploymentStatus = append(deploymentStatus, statusInfo)
	}

	c.JSON(200, gin.H{
		"Status": deploymentStatus,
	})
}
