package global

import (
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/listers/apps/v1"
)

var (
	SharedInformers informers.SharedInformerFactory // SharedInformers 全局共享的 informers
	Deployments     v1.DeploymentLister             // 部署列表器
	Clientset       kubernetes.Interface            // Kubernetes 客户端
)

//var SharedInformers informers.SharedInformerFactory                  // SharedInformers 全局共享的 informers
//var Deployments = SharedInformers.Apps().V1().Deployments().Lister() // 建一个 Deployments 用于列出所有的 Deployments
