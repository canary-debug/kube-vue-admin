package get_nodes_resources

import (
	"net/http"

	"github.com/canary-debug/kube-vue-admin/pkg/global"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func GetNodeLen(c *gin.Context) {
	// 使用 global.Nodes.Lister() 获取所有节点，然后计算长度
	nodes, err := global.Nodes.List(labels.Everything())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get nodes: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"node_len": len(nodes),
	})
}

func GetNodes(c *gin.Context) {
	//// 获取所有节点信息
	//node, err := resources.Clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{
	//		"error": "Failed to get nodes: " + err.Error(),
	//	})
	//	return
	//}
	//
	//var response []map[string]interface{}
	//for _, n := range node.Items {
	//	// 污点字段
	//	var taints []string
	//	for _, taint := range n.Spec.Taints {
	//		taints = append(taints, taint.Key+"="+taint.Value+":"+string(taint.Effect))
	//	}
	//
	//	// 根据条件确定节点状态
	//	status := "Unknown"
	//	for _, condition := range n.Status.Conditions {
	//		if condition.Type == "Ready" {
	//			if condition.Status == "True" {
	//				status = "Ready"
	//			} else {
	//				status = "NotReady"
	//			}
	//			break
	//		}
	//	}
	//
	//	// 获取节点角色
	//	role := getNodeRole(n)
	//
	//	response = append(response, map[string]interface{}{
	//		"name":   n.Name,
	//		"status": status,
	//		"taints": taints,
	//		"role":   role,
	//		"cpu":    n.Status.Capacity.Cpu().String(),
	//		"memory": n.Status.Capacity.Memory().String(),
	//	})
	//}
	//
	//c.JSON(http.StatusOK, response)
	// 1. 安全检查：确保 Lister 已经初始化
	if global.Nodes == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Nodes Informer 尚未初始化"})
		return
	}

	// 2. 从本地缓存获取所有节点
	nodes, err := global.Nodes.List(labels.Everything())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取缓存数据失败: " + err.Error()})
		return
	}

	var response []map[string]interface{}
	for _, n := range nodes {
		// 污点字段处理
		var taints []string
		for _, taint := range n.Spec.Taints {
			taints = append(taints, taint.Key+"="+taint.Value+":"+string(taint.Effect))
		}

		// 根据状态条件确定节点是否 Ready
		status := "Unknown"
		for _, condition := range n.Status.Conditions {
			if condition.Type == v1.NodeReady { // 使用官方常量更规范
				if condition.Status == v1.ConditionTrue {
					status = "Ready"
				} else {
					status = "NotReady"
				}
				break
			}
		}

		// 获取节点角色 (调用下方优化后的函数)
		role := getNodeRole(*n)

		response = append(response, map[string]interface{}{
			"name":   n.Name,
			"status": status,
			"taints": taints,
			"role":   role,
			"cpu":    n.Status.Capacity.Cpu().String(),
			"memory": n.Status.Capacity.Memory().String(),
			"labels": n.Labels, // 顺便返回标签，方便前端展示
		})
	}

	c.JSON(http.StatusOK, response)
}

// getNodeRole 获取节点角色
func getNodeRole(node v1.Node) string {
	if node.Labels == nil {
		return "worker"
	}

	// 常见的 Master/Control-Plane 角色标签
	masterLabels := []string{
		"node-role.kubernetes.io/master",
		"node-role.kubernetes.io/control-plane",
	}

	for _, label := range masterLabels {
		if _, has := node.Labels[label]; has {
			return "master"
		}
	}

	return "worker"
}
