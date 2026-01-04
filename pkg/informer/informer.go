package informer

import (
	"context"
	"log"

	"github.com/canary-debug/kube-vue-admin/api/resources"
	"github.com/canary-debug/kube-vue-admin/pkg/global"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/informers"
)

func Informer() {
	// 创建一个上下文，用于取消请求
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 延迟调用取消请求

	// 获取所有 Namespace 的 SharedInformer
	SharedInformers := informers.NewSharedInformerFactory(resources.Clientset, 0)

	// 声明获取资源的 GroupVersionResource
	gvrs := []schema.GroupVersionResource{
		{Group: "", Version: "v1", Resource: "pods"},
		{Group: "apps", Version: "v1", Resource: "deployments"},
	}

	// 为每个 GroupVersionResource 创建 SharedInformer
	for _, gvr := range gvrs {
		if _, err := SharedInformers.ForResource(gvr); err != nil {
			panic(err)
		}
	}

	// 启动所有 SharedInformer
	SharedInformers.Start(ctx.Done())

	// 等待所有 SharedInformer 同步
	SharedInformers.WaitForCacheSync(ctx.Done())

	// 关键修改：设置全局变量
	global.SharedInformers = SharedInformers
	global.Deployments = SharedInformers.Apps().V1().Deployments().Lister()

	// 打印一下信息
	log.Println("所有 Informer 启动成功！")

}
