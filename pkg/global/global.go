package global

import (
	"k8s.io/client-go/informers"
	v1 "k8s.io/client-go/listers/apps/v1"
	corev1 "k8s.io/client-go/listers/core/v1"
)

var (
	SharedInformers informers.SharedInformerFactory // SharedInformers 全局共享的 informers
	Deployments     v1.DeploymentLister             // 部署列表器
	Daemonsets      v1.DaemonSetLister              // 守护进程列表器
	Statefulset     v1.StatefulSetLister
	Nodes           corev1.NodeLister      // 节点列表器
	Pods            corev1.PodLister       // Pod列表器
	Namespaces      corev1.NamespaceLister // 命名空间列表器
)
