package api

import (
	"errors"
	"log"

	"github.com/canary-debug/kube-vue-admin/database"
	"github.com/canary-debug/kube-vue-admin/models"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// 用户登录函数
func Login(c *gin.Context) {
	// 用户名和密码参数
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 从库里面查用户名和密码
	var user models.Users

	// result 是一个数据库查询结果
	result := database.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {

		// 特别检查：是否是“记录未找到”错误 (gorm.ErrRecordNotFound)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("查询成功，但未找到用户: %s", username)
			// 此时 user 变量会是零值 (e.g., ID=0, Username="")
			return
		}

		// 处理其他严重的数据库错误（如连接断开、SQL语法错误等）
		log.Fatalf("❌ 数据库查询失败: %v", result.Error)
		return
	}

	// 参数校验
	if username == "" || password == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "用户名和密码不能为空",
		})
		return
	}

	// 验证密码
	if user.Password != password {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "用户名或密码错误",
		})
		return
	}

	// 登录成功
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{
			"user_id":  user.ID,
			"username": user.Username,
			"password": user.Password,
		},
	})

}
