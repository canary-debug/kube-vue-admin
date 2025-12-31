package get_nodes_resources

import (
	"context"
	"net/http"

	"github.com/canary-debug/kube-vue-admin/api/resources"
	"github.com/gin-gonic/gin"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetNodeName(c *gin.Context) {

	// 获取所有节点信息
	node, err := resources.Clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get nodes: " + err.Error(),
		})
		return
	}

	// 构建每个节点详细信息
	var response []map[string]interface{}
	for _, n := range node.Items {

		// 根据条件确定节点状态
		status := "Unknown"
		for _, condition := range n.Status.Conditions {
			if condition.Type == "Ready" {
				if condition.Status == "True" {
					status = "Ready"
				} else {
					status = "NotReady"
				}
				break
			}
		}

		// 污点字段
		var taints []string
		for _, taint := range n.Spec.Taints {
			taints = append(taints, taint.Key+"="+taint.Value+":"+string(taint.Effect))
		}

		// 获取节点IP地址
		IP := getNodeIP(n)

		response = append(response, map[string]interface{}{
			"name":           n.Name,
			"status":         status,
			"role":           getNodeRole(n), // 获取节点角色
			"ip":             IP,
			"osType":         n.Status.NodeInfo.OperatingSystem,         // 操作系统类型
			"osVersion":      n.Status.NodeInfo.OSImage,                 // 操作系统版本
			"kubeletVersion": n.Status.NodeInfo.KubeletVersion,          // kubelet 版本
			"kube-proxy":     n.Status.NodeInfo.KubeProxyVersion,        // kube-proxy 版本
			"dockerVersion":  n.Status.NodeInfo.ContainerRuntimeVersion, // docker 运行时版本
			"coreVersion":    n.Status.NodeInfo.KernelVersion,           // 内核版本
			"nodecreatetime": n.ObjectMeta.CreationTimestamp.String(),   // 节点创建时间
			"taints":         taints,
		})
	}

	c.JSON(http.StatusOK, response)
}

// 获取节点IP
func getNodeIP(n v1.Node) interface{} {
	// 获取节点IP
	for _, address := range n.Status.Addresses {
		if address.Type == "InternalIP" {
			return address.Address
		}
	}
	return "N/A"

}
