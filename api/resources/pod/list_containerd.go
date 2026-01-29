package pod

import (
	"context"
	"net/http"

	"github.com/canary-debug/kube-vue-admin/api/resources"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetPodContainersHandler 直接作为 Gin 的 Handler
func GetPodContainersHandler(c *gin.Context) {
	// 直接从 URL 路径参数获取 (例如: /:namespace/:podname/containers)
	namespace := c.Param("namespace")
	podName := c.Param("podname")

	// 1. 调用你已经初始化好的 resources.Clientset (假设变量名为 GlobalClient)
	// 如果命名空间或 Pod 名字为空，Get 内部会报错，我们直接处理 err
	pod, err := resources.Clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "找不到指定的 Pod",
			"error":   err.Error(),
		})
		return
	}

	// 2. 准备容器状态映射表 (Status)
	statusMap := make(map[string]struct {
		restartCount int32
		state        string
	})
	for _, s := range pod.Status.ContainerStatuses {
		var stateStr string
		switch {
		case s.State.Running != nil:
			stateStr = "Running"
		case s.State.Waiting != nil:
			stateStr = "Waiting"
		case s.State.Terminated != nil:
			stateStr = "Terminated"
		default:
			stateStr = "Unknown"
		}
		statusMap[s.Name] = struct {
			restartCount int32
			state        string
		}{restartCount: s.RestartCount, state: stateStr}
	}

	// 3. 构造返回数组 (结合 Spec 和 Status)
	type containerInfo struct {
		Name         string  `json:"name"`
		RestartCount int32   `json:"restart_count"`
		State        string  `json:"state"`
		Ports        []int32 `json:"ports"`
	}

	results := make([]containerInfo, 0, len(pod.Spec.Containers))

	for _, container := range pod.Spec.Containers {
		ports := []int32{}
		for _, p := range container.Ports {
			ports = append(ports, p.ContainerPort)
		}

		info := containerInfo{
			Name:  container.Name,
			Ports: ports,
		}

		// 填充来自 Status 的动态数据
		if stat, ok := statusMap[container.Name]; ok {
			info.RestartCount = stat.restartCount
			info.State = stat.state
		}

		results = append(results, info)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": results,
	})
}
