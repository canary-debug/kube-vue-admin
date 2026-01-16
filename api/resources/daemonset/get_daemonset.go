package daemonset

import (
	"log"

	"github.com/canary-debug/kube-vue-admin/pkg/global"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/labels"
)

func GetDaemonset(c *gin.Context) {
	// 创建一个切片，用于存储 deployment 的名字和状态信息
	//deploymentStatus := []interface{}{}
	namespace := c.Param("namespace")
	log.Println("namespace: ", namespace)

	// 使用 informer 的 lister 来获取指定命名空间下的 daemonset 列表
	daemonsetLister := global.Daemonsets

	// 通过 namespace filter 获取 deployments
	daemonsets, err := daemonsetLister.DaemonSets(namespace).List(labels.Everything())
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	// 遍历 deployments 返回
	for _, daemonset := range daemonsets {
		// 获取 deployment 的状态信息
		statusInfo := map[string]interface{}{
			"name":     daemonset.Name,                          // daemonset 的名字
			"replicas": daemonset.Status.DesiredNumberScheduled, // daemonset 的副本数
		}
		c.JSON(200, gin.H{
			"status": statusInfo,
		})

		// 获取运行状态
		//var runningStatus string
		//for _, condition := range daemonset.Status.Conditions {
		//	if condition.Type == "Available" {
		//		if condition.Status == "True" {
		//			// 使用格式化字符串构建 (实际/期望) 格式
		//			runningStatus = fmt.Sprintf("运行中 (%d/%d)",
		//				daemonset.Status.ReadyReplicas,
		//				*daemonset.Spec.Replicas)
		//		} else {
		//			runningStatus = fmt.Sprintf("未就绪 (%d/%d)",
		//				daemonset.Status.ReadyReplicas,
		//				*daemonset.Spec.Replicas)
		//		}
		//		break
		//	}
		//}
		//statusInfo["status"] = runningStatus

	}
}
