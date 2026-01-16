package global

import (
	"k8s.io/client-go/informers"
	v1 "k8s.io/client-go/listers/apps/v1"
	corev1 "k8s.io/client-go/listers/core/v1"
)

var (
	SharedInformers informers.SharedInformerFactory // SharedInformers 全局共享的 informers
	Deployments     v1.DeploymentLister             // 部署列表器
	Daemonsets      v1.DaemonSetLister
	Nodes           corev1.NodeLister
	Pods            corev1.PodLister
	Namespaces      corev1.NamespaceLister
)
