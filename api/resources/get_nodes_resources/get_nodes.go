package get_nodes_resources

import (
	"context"
	"net/http"

	"github.com/canary-debug/kube-vue-admin/api/resources"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetNodes(c *gin.Context) {
	// 获取所有节点信息
	node, err := resources.Clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get nodes: " + err.Error(),
		})
		return
	}

	var response []map[string]interface{}
	for _, n := range node.Items {
		// 污点字段
		var taints []string
		for _, taint := range n.Spec.Taints {
			taints = append(taints, taint.Key+"="+taint.Value+":"+string(taint.Effect))
		}

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

		// 获取节点角色
		role := getNodeRole(n)

		response = append(response, map[string]interface{}{
			"name":   n.Name,
			"status": status,
			"taints": taints,
			"role":   role,
			"cpu":    n.Status.Capacity.Cpu().String(),
			"memory": n.Status.Capacity.Memory().String(),
		})
	}

	c.JSON(http.StatusOK, response)
}

// getNodeRole 获取节点角色
func getNodeRole(node v1.Node) string {
	labels := node.Labels

	// 检查master标签
	if _, hasMasterLabel := labels["node-role.kubernetes.io/master"]; hasMasterLabel {
		return "master"
	}

	// 检查control-plane标签
	if _, hasControlPlaneLabel := labels["node-role.kubernetes.io/control-plane"]; hasControlPlaneLabel {
		return "master"
	}

	// 默认为worker节点
	return "worker"
}
