package resources

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPod(c *gin.Context) {
	// 获取命名空间, 资源名称 返回 pod 信息
	namespace := c.Param("namespace")
	resourcename := c.Param("resourcename")
	if namespace == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "namespace parameter is required",
		})
		return
	} else if resourcename == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "resourcename parameter is required",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"resourcename": resourcename,
		"namespace":    namespace,
	})

}
