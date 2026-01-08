package deployment

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/canary-debug/kube-vue-admin/api/resources"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RestartDeploymentRequest 重启请求结构体
type RestartDeploymentRequest struct {
	Name      string `json:"name" binding:"required"`
	Namespace string `json:"namespace" binding:"required"`
}

// RestartDeployment 重启Deployment接口
func RestartDeployment(c *gin.Context) {
	var req RestartDeploymentRequest

	// 绑定JSON请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("参数绑定失败: %v", err),
		})
		return
	}

	// 验证ClientSet是否已初始化
	if resources.Clientset == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Kubernetes客户端未初始化",
		})
		return
	}

	// 执行重启操作
	err := restartDeploymentFunc(req.Name, req.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Deployment %s 在命名空间 %s 中已触发重启", req.Name, req.Namespace),
	})
}

// restartDeploymentFunc 实际执行重启逻辑的函数
func restartDeploymentFunc(name, namespace string) error {
	ctx := context.Background()

	// 获取当前Deployment
	deployment, err := resources.Clientset.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return fmt.Errorf("deployment %s 在命名空间 %s 中不存在", name, namespace)
		}
		return fmt.Errorf("获取Deployment失败: %v", err)
	}

	// 在Pod模板中添加重启时间戳注解，触发滚动更新
	if deployment.Spec.Template.ObjectMeta.Annotations == nil {
		deployment.Spec.Template.ObjectMeta.Annotations = make(map[string]string)
	}
	deployment.Spec.Template.ObjectMeta.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

	// 更新Deployment以触发重启
	_, err = resources.Clientset.AppsV1().Deployments(namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("更新Deployment失败: %v", err)
	}
	log.Println("Deployment %s 在命名空间 %s 中已触发重启", name, namespace)

	return nil
}
