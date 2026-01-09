package etcd

import (
	"net/http"
	"strings"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection" // 辅助构造选择器

	"github.com/canary-debug/kube-vue-admin/pkg/global"
	"github.com/gin-gonic/gin"
)

// EtcdStatus 结构体保持不变，兼容前端
type EtcdStatus struct {
	Name    string `json:"name"`
	Healthy bool   `json:"healthy"`
	Message string `json:"message,omitempty"`
	IP      string `json:"ip,omitempty"` // 新增：展示 Pod IP 更有用
}

func GetEtcdStatus(c *gin.Context) {
	// 1. 验证 Pod Informer 是否就绪
	if global.Pods == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Pod Informer 未初始化"})
		return
	}

	// 2. 构造选择器：通常 etcd 的 Pod 带有标签 component=etcd
	// 也可以简单用 labels.Everything() 然后在循环里用 strings.Contains 判断名称
	selector := labels.NewSelector()
	requirement, _ := labels.NewRequirement("component", selection.Equals, []string{"etcd"})
	selector = selector.Add(*requirement)

	// 3. 在 kube-system 命名空间下列出 Pod
	pods, err := global.Pods.Pods("kube-system").List(selector)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取 Pod 失败: " + err.Error()})
		return
	}

	var etcdStatuses []EtcdStatus
	for _, pod := range pods {
		status := EtcdStatus{
			Name:    pod.Name,
			Healthy: false,
			IP:      pod.Status.PodIP,
		}

		// 4. 判断健康状态：Pod 处于 Running 阶段且容器全部 Ready
		if pod.Status.Phase == v1.PodRunning {
			allReady := true
			for _, cs := range pod.Status.ContainerStatuses {
				if !cs.Ready {
					allReady = false
					status.Message = cs.State.Waiting.Message // 记录错误原因
					break
				}
			}
			status.Healthy = allReady
		} else {
			status.Message = "Pod 状态: " + string(pod.Status.Phase)
		}

		etcdStatuses = append(etcdStatuses, status)
	}

	// 如果没搜到标签，可能是某些集群标签不同，尝试按名字模糊匹配
	if len(etcdStatuses) == 0 {
		allPods, _ := global.Pods.Pods("kube-system").List(labels.Everything())
		for _, p := range allPods {
			if strings.Contains(p.Name, "etcd") {
				// 执行同样的健康检查逻辑 (建议抽离成函数)
				etcdStatuses = append(etcdStatuses, convertPodToEtcdStatus(p))
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    etcdStatuses,
		"total":   len(etcdStatuses),
		"message": "通过 Pod 状态获取 etcd 信息成功",
	})
}

func convertPodToEtcdStatus(pod *v1.Pod) EtcdStatus {
	healthy := pod.Status.Phase == v1.PodRunning
	msg := ""
	if !healthy {
		msg = "Pod is " + string(pod.Status.Phase)
	}
	return EtcdStatus{
		Name:    pod.Name,
		Healthy: healthy,
		Message: msg,
		IP:      pod.Status.PodIP,
	}
}
