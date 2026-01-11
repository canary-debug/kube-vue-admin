package main

import (
	config "github.com/canary-debug/kube-vue-admin/configs"
	"io"
	"os"

	"github.com/canary-debug/kube-vue-admin/api"
	"github.com/canary-debug/kube-vue-admin/api/resources"
	"github.com/canary-debug/kube-vue-admin/api/resources/deployment"
	"github.com/canary-debug/kube-vue-admin/api/resources/etcd"
	"github.com/canary-debug/kube-vue-admin/api/resources/get_nodes_resources"
	"github.com/canary-debug/kube-vue-admin/database"
	"github.com/canary-debug/kube-vue-admin/pkg/informer"
	"github.com/canary-debug/kube-vue-admin/tokens"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 生产模式发布
	ginMode := os.Getenv("GIN_MODE")
	switch ginMode {
	case "release":
		gin.SetMode(gin.ReleaseMode) // 生产环境：禁用调试信息，性能更优
	case "debug":
		gin.SetMode(gin.DebugMode) // 可选：手动指定调试模式
	default:
		gin.SetMode(gin.DebugMode) // 默认：开发环境（dev分支）使用调试模式
	}

	// 初始化配置
	if err := config.Init(); err != nil {
		panic("配置初始化失败: " + err.Error())
	}

	// 加载 informer
	// 在一个独立的协程里启动，不要让它阻塞主流程，但要保证 stopCh 没关闭
	stopCh := make(chan struct{})
	go informer.StartInformer(stopCh)

	// 加载密钥
	tokens.LoadToekn()

	// 生成路由
	r := gin.Default()

	// 全局启用 CORS 中间件（推荐）
	r.Use(cors())

	/*
		创建记录日志的文件
		将日志同时写入文件和控制台
	*/
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	// 添加日志和恢复中间件
	r.Use(gin.Logger()) // 这会记录所有请求日志
	r.Use(gin.Recovery())

	// 初始化 Mysql
	database.ConnectDatabase()
	// 初始化 Redis
	database.ConnectRedis()

	// 登陆注册路由
	authRoutes := r.Group("/api/auth")
	{
		// 登录接口
		authRoutes.POST("/login", api.Login)

		//  注册接口
		authRoutes.POST("/register", api.Register)

	}

	// 退出接口 todo (还没实现退出是返回退出的用户名字)
	privateAuth := r.Group("/api/auth")
	privateAuth.Use(tokens.JWTMiddleware()) // 给这个组应用中间件
	{
		privateAuth.POST("/logout", api.LoginOut)
	}

	// 受保护的接口, 需要带上 Token
	// 控制器资源路由
	k8sRoutes := r.Group("/api/k8s", tokens.JWTMiddleware())
	{

		// GetControllersInNamespace 获取指定命名空间下的所有控制器资源
		k8sRoutes.GET("/namespaces/:namespace/controllers", resources.GetControllersInNamespace)
		// 获取命名空间, 资源名称 返回 pod 信息
		k8sRoutes.GET("/:namespace/:resourcename/pods", resources.GetPod)
		// 获取节点资源
		k8sRoutes.GET("/get/nodes", get_nodes_resources.GetNodes)
		// 获取节点个数
		k8sRoutes.GET("/get/nodes/len", get_nodes_resources.GetNodeLen)
		// 获取节点详细信息
		k8sRoutes.GET("/get/nodename", get_nodes_resources.GetNodeName)
		// 获取容器组
		k8sRoutes.GET("/get/container_group/:namespace", get_nodes_resources.GetContainerGroup)
		// 获取 namespaces 的个数
		k8sRoutes.GET("/get/namespaces", get_nodes_resources.GetNameSpacesLen)
		// 获取所有命名空间的名字
		k8sRoutes.GET("/get/namespaces/namespacename", get_nodes_resources.GetNameSpaces)
		// 获取指定 deployment
		k8sRoutes.GET("/get/deployment/:namespace", get_nodes_resources.GetDeployment)
		// 获取 pod 个数
		k8sRoutes.GET("/get/pods/len", resources.GetPodCount)
		// 获取集群是否健康
		k8sRoutes.GET("/get/cluster_healthz", get_nodes_resources.Get_Cluster_Healthz)
		// 重启 deployment
		k8sRoutes.POST("/restart/deployment", deployment.RestartDeployment)
		// 获取 deployment 下的所有 pods
		k8sRoutes.POST("/deployment/pods", deployment.Get_Deployment_Pods)
		// 获取etcd状态
		k8sRoutes.GET("/etcd/status", etcd.GetEtcdStatus)
		// 获取 pod 日志
		k8sRoutes.GET("/pod/logs/:namespace/:pod", deployment.GetPodLogs)
		// 天气接口
		k8sRoutes.GET("/weather", api.Weather)
	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// listen 9000 prot
	if err := r.Run(":9000"); err != nil {
		panic(err)
	}
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*") // 或指定域名，如 "http://your-frontend.com"
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204) // 预检成功，直接返回 204
			return
		}
		c.Next()
	}
}
