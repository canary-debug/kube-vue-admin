package service

import (
	"context"
	"log"
	"net/http"

	"github.com/canary-debug/kube-vue-admin/api/resources"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeleteService(c *gin.Context) {
	namespace := c.Param("namespace")
	servicename := c.Param("servicename")

	log.Printf("删除 Service - 命名空间: %s, Service名称: %s", namespace, servicename)

	err := resources.Clientset.CoreV1().Services(namespace).Delete(context.TODO(), servicename, metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Service 不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Service %s 删除成功", servicename)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Service " + servicename + " 删除成功",
	})
}
