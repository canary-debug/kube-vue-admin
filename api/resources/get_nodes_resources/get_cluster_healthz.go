package get_nodes_resources

import (
	"context"

	"github.com/canary-debug/kube-vue-admin/api/resources"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//func Get_Cluster_Healthz(c *gin.Context) {
//	// 使用 DiscoveryClient 或 RestClient 访问健康检查接口
//	content, err := resources.Clientset.Discovery().RESTClient().Get().AbsPath("/readyz").DoRaw(context.TODO())
//	if err != nil {
//		fmt.Println("Error getting cluster healthz:", err)
//		c.JSON(500, gin.H{"error": false, "message": err.Error()})
//		return // 添加 return 语句
//	}
//
//	// 检查返回内容是否表示健康状态
//	response := string(content)
//	if response != "ok" {
//		c.JSON(500, gin.H{"error": false, "status": response})
//		return
//	}
//
//	c.JSON(200, gin.H{"status": response})
//}

func Get_Cluster_Healthz(c *gin.Context) {
	// 假设你已经在 resources 里初始化了 PodInformer
	// 如果没有，直接用 Clientset List 也可以
	ctx := context.TODO()

	// 1. 获取所有非 Ready 的 Pod
	pods, _ := resources.Clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})

	// 统计非 Ready 的 Pod 数量
	unhealthyCount := 0
	crashLoopCount := 0

	for _, pod := range pods.Items {
		// 检查重启次数过多的 Pod (比如最近 10 分钟重启超过 3 次)
		for _, status := range pod.Status.ContainerStatuses {
			if status.RestartCount > 3 {
				crashLoopCount++
			}
			if !status.Ready {
				unhealthyCount++
			}
		}
	}

	// 2. 检查关键 Endpoints (比如 CoreDNS)
	// 如果删除了 Calico，这里的 Subsets 会是空的
	ep, _ := resources.Clientset.CoreV1().Endpoints("kube-system").Get(ctx, "kube-dns", metav1.GetOptions{})
	networkFunctional := len(ep.Subsets) > 0

	// 3. 综合评估
	if !networkFunctional || crashLoopCount > 5 {
		c.JSON(500, gin.H{
			//"status": "Critical",
			"status": "Not Health",
			"reason": "Network plugin failure or excessive pod crashes",
			"details": gin.H{
				"unhealthy_pods": unhealthyCount,
				"crash_pods":     crashLoopCount,
				"dns_active":     networkFunctional,
			},
		})
		return
	}

	c.JSON(200, gin.H{"status": "Healthy"})
}
