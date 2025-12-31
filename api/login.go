package api

import (
	"errors"
	"log"
	"net/http"

	"github.com/canary-debug/kube-vue-admin/database"
	"github.com/canary-debug/kube-vue-admin/models"
	"github.com/canary-debug/kube-vue-admin/tokens"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginRequest struct {
	ID       uint   `json:"id"`
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6,max=30"`
}

// 用户登录函数
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 从库里面查用户名和密码
	var user models.Users

	// result 是一个数据库查询结果
	result := database.DB.Where("username = ?", req.Username).First(&user)
	if result.Error != nil {

		// 特别检查：是否是“记录未找到”错误 (gorm.ErrRecordNotFound)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("查询成功，但未找到用户: %s", req.Username)
			c.JSON(200, gin.H{
				"code": 400,
				"msg":  "用户不存在",
			})
			// 此时 user 变量会是零值 (e.g., ID=0, Username="")
			return
		}

		// 处理其他严重的数据库错误（如连接断开、SQL语法错误等）
		log.Fatalf("❌ 数据库查询失败: %v", result.Error)
		return
	}

	// 参数校验
	if req.Username == "" || req.Password == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "用户名和密码不能为空",
		})
		return
	}

	// 比对密码：传入【用户输入的明文密码】和【数据库存储的哈希密码】
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		// 密码不匹配（bcrypt.ErrMismatchedHashAndPassword）
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成 Token
	token, err := tokens.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
	}

	// 登录成功, 返回 Toekn
	c.JSON(200, gin.H{
		"code":  200,
		"token": token,
	})

}
