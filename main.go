package main

import (
	"io"
	"os"

	"github.com/canary-debug/kube-vue-admin/api"
	"github.com/canary-debug/kube-vue-admin/database"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func init() {

}

func main() {
	// 生成路由
	r := gin.Default()

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

	// 路由分组
	authRoutes := r.Group("/api/auth")
	{
		// 登录接口
		authRoutes.POST("/login", api.Login)

		//  注册接口
		authRoutes.POST("/register", api.Register)
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
