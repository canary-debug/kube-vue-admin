package daemonset

import (
	"fmt"
	"log"

	"github.com/canary-debug/kube-vue-admin/pkg/global"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/labels"
)

func GetDaemonset(c *gin.Context) {
	// 创建一个切片，用于存储 deployment 的名字和状态信息
	daemonSetStatus := []interface{}{}
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

	// 遍历 daemonsets 返回
	for _, ds := range daemonsets {
		statusInfo := map[string]interface{}{
			"name": ds.Name,
		}

		// 获取运行状态 (DaemonSet 特有的状态字段)
		// DesiredNumberScheduled: 应该调度的节点数
		// NumberReady: 已就绪的节点数
		desired := ds.Status.DesiredNumberScheduled
		ready := ds.Status.NumberReady

		var runningStatus string
		if ready == desired && desired > 0 {
			runningStatus = fmt.Sprintf("运行中 (%d/%d)", ready, desired)
		} else {
			runningStatus = fmt.Sprintf("同步中/部分就绪 (%d/%d)", ready, desired)
		}

		statusInfo["status"] = runningStatus
		//statusInfo["desired"] = desired
		statusInfo["replicas"] = ready

		// 5. 从 ManagedFields 获取最后更新时间
		if len(ds.ManagedFields) > 0 {
			lastField := ds.ManagedFields[len(ds.ManagedFields)-1]
			if lastField.Time != nil {
				formattedTime := lastField.Time.Format("2006-01-02 15:04:05")
				statusInfo["update_time"] = formattedTime
			}
		}

		daemonSetStatus = append(daemonSetStatus, statusInfo)
	}

	c.JSON(200, gin.H{
		"Status": daemonSetStatus,
	})

}
