package main

import (
	"github.com/canary-debug/kube-vue-admin/api"
	"github.com/canary-debug/kube-vue-admin/database"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"os"
)

func main() {
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

	// 初始化数据库连接和表
	database.ConnectDatabase()

	// 登陆注册路由
	authRoutes := r.Group("/api/auth")
	{
		// 登录接口
		authRoutes.POST("/login", api.Login)

		//  注册接口
		authRoutes.POST("/register", api.Register)
	}

	// 控制器资源路由
	k8sRoutes := r.Group("/api/k8s")
	{
		k8sRoutes.GET("/namespaces/:namespace/controllers", api.GetControllersInNamespace)
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
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204) // 预检成功，直接返回 204
			return
		}
		c.Next()
	}
}
