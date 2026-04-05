package service

import (
	"log"
	"net/http"

	"github.com/canary-debug/kube-vue-admin/pkg/global"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func GetServices(c *gin.Context) {
	namespace := c.Param("namespace")
	log.Println("获取 Service - 命名空间: ", namespace)

	if global.Services == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Services Informer 尚未初始化完成，请稍后重试"})
		return
	}

	var services []interface{}

	if namespace == "" || namespace == "all" {
		namespaces, err := global.Namespaces.List(labels.Everything())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取命名空间失败: " + err.Error()})
			return
		}

		for _, ns := range namespaces {
			nsServices, err := global.Services.Services(ns.Name).List(labels.Everything())
			if err != nil {
				log.Printf("获取命名空间 %s 的 Service 失败: %v", ns.Name, err)
				continue
			}

			for _, svc := range nsServices {
				services = append(services, formatServiceInfo(svc))
			}
		}
	} else {
		nsServices, err := global.Services.Services(namespace).List(labels.Everything())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取 Service 失败: " + err.Error()})
			return
		}

		for _, svc := range nsServices {
			services = append(services, formatServiceInfo(svc))
		}
	}

	log.Printf("共获取到 %d 个 Service", len(services))

	c.JSON(http.StatusOK, gin.H{
		"services": services,
		"total":    len(services),
	})
}

func formatServiceInfo(svc *v1.Service) map[string]interface{} {
	serviceInfo := make(map[string]interface{})

	serviceInfo["name"] = svc.GetName()
	serviceInfo["namespace"] = svc.GetNamespace()

	if len(svc.Spec.ExternalIPs) > 0 {
		serviceInfo["external_ip"] = svc.Spec.ExternalIPs[0]
	} else {
		serviceInfo["external_ip"] = nil
	}

	if len(svc.Spec.Ports) > 0 {
		for _, port := range svc.Spec.Ports {
			if port.NodePort != 0 {
				serviceInfo["port"] = port.Port
				serviceInfo["node_port"] = port.NodePort
				serviceInfo["target_port"] = port.TargetPort.String()
				serviceInfo["protocol"] = string(port.Protocol)
				break
			}
		}
	}

	if _, exists := serviceInfo["port"]; !exists {
		serviceInfo["port"] = nil
		serviceInfo["node_port"] = nil
		serviceInfo["target_port"] = nil
		serviceInfo["protocol"] = nil
	}

	if len(svc.Spec.ClusterIPs) > 0 {
		serviceInfo["cluster_ip"] = svc.Spec.ClusterIPs[0]
	} else {
		serviceInfo["cluster_ip"] = nil
	}

	serviceInfo["creation_timestamp"] = svc.GetCreationTimestamp().Format("2006-01-02 15:04:05")

	return serviceInfo
}
