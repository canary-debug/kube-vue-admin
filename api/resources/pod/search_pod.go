package pod

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/canary-debug/kube-vue-admin/pkg/global"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/labels"
)

// SearchPodHandler 搜索 Pod 的处理函数
func SearchPodHandler(c *gin.Context) {
	// 从查询参数获取 filter (正则表达式)
	filter := c.Query("filter")

	if filter == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "filter parameter is required",
		})
		return
	}

	// 预编译正则表达式以检查有效性
	reg, err := regexp.Compile(filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid regex pattern: %v", err),
		})
		return
	}

	// 验证全局 Informer 是否已初始化
	if global.Pods == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "global Pod informer not initialized",
		})
		return
	}

	// 始终搜索所有命名空间
	pods, err := global.Pods.List(labels.Everything())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to list pods: %v", err),
		})
		return
	}

	// 进行正则匹配过滤
	var matchedPodInfos []PodInfo
	for _, p := range pods {
		if reg.MatchString(p.Name) {
			podInfo := ConvertPodToPodInfo(p)
			matchedPodInfos = append(matchedPodInfos, *podInfo)
		}
	}

	// 返回匹配结果
	c.JSON(http.StatusOK, gin.H{
		"data":    matchedPodInfos,
		"total":   len(matchedPodInfos),
		"message": "success",
	})
}
