package informer

import (
	"log"

	"github.com/canary-debug/kube-vue-admin/api/resources"
	"github.com/canary-debug/kube-vue-admin/pkg/global"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/informers"
)

func StartInformer(stopCh <-chan struct{}) {
	// 1. 使用全局 Clientset 创建 Factory
	SharedInformers := informers.NewSharedInformerFactory(resources.Clientset, 0)

	// 2. 注册资源（你的 gvrs 逻辑没问题）
	gvrs := []schema.GroupVersionResource{
		{Group: "", Version: "v1", Resource: "nodes"},
		{Group: "", Version: "v1", Resource: "pods"},
		{Group: "apps", Version: "v1", Resource: "deployments"},
		{Group: "", Version: "v1", Resource: "namespaces"},
	}
	for _, gvr := range gvrs {
		SharedInformers.ForResource(gvr)
	}

	// 3. 启动并等待同步。注意这里传入的是 stopCh
	SharedInformers.Start(stopCh)
	SharedInformers.WaitForCacheSync(stopCh)

	// 4. 赋值全局变量
	global.SharedInformers = SharedInformers
	global.Deployments = SharedInformers.Apps().V1().Deployments().Lister()
	global.Nodes = SharedInformers.Core().V1().Nodes().Lister()
	global.Pods = SharedInformers.Core().V1().Pods().Lister()
	global.Namespaces = SharedInformers.Core().V1().Namespaces().Lister()

	log.Println("所有 Informer 启动同步完成，正在持续监听...")
}
