package statefulset

import (
	"fmt"
	"log"

	"github.com/canary-debug/kube-vue-admin/pkg/global"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/labels"
)

func GetStatefulset(c *gin.Context) {
	// 创建一个切片，用于存储 statefulset 的名字和状态信息
	statefulSetStatus := []interface{}{}
	namespace := c.Param("namespace")
	log.Println("namespace: ", namespace)

	// 1. 使用 informer 的 lister 来获取指定命名空间下的 statefulset 列表
	// 确保 global.Statefulsets 在初始化时已经注册
	statefulsetLister := global.Statefulset
	if statefulsetLister == nil {
		c.JSON(500, gin.H{
			"error": "StatefulSet Informer 未初始化",
		})
		return
	}

	// 2. 通过 namespace 筛选获取 statefulsets
	statefulsets, err := statefulsetLister.StatefulSets(namespace).List(labels.Everything())
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 3. 遍历 statefulsets
	for _, sts := range statefulsets {
		statusInfo := map[string]interface{}{
			"name": sts.Name,
		}

		// 4. 获取运行状态 (StatefulSet 特有的状态字段)
		// Replicas: 期望的副本数
		// ReadyReplicas: 当前已就绪的副本数
		desired := sts.Spec.Replicas
		var desiredVal int32
		if desired != nil {
			desiredVal = *desired
		}
		ready := sts.Status.ReadyReplicas

		var runningStatus string
		if ready == desiredVal && desiredVal > 0 {
			runningStatus = fmt.Sprintf("运行中 (%d/%d)", ready, desiredVal)
		} else if desiredVal == 0 {
			runningStatus = "已停止 (0/0)"
		} else {
			runningStatus = fmt.Sprintf("同步中 (%d/%d)", ready, desiredVal)
		}

		statusInfo["status"] = runningStatus
		statusInfo["replicas"] = ready

		// 5. 从 ManagedFields 获取最后更新时间
		if len(sts.ManagedFields) > 0 {
			lastField := sts.ManagedFields[len(sts.ManagedFields)-1]
			if lastField.Time != nil {
				formattedTime := lastField.Time.Format("2006-01-02 15:04:05")
				statusInfo["update_time"] = formattedTime
			}
		}

		statefulSetStatus = append(statefulSetStatus, statusInfo)
	}

	c.JSON(200, gin.H{
		"Status": statefulSetStatus,
	})
}
