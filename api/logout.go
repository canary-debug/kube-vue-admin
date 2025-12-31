package api

import (
	"log"
	"net/http"

	"github.com/canary-debug/kube-vue-admin/database"
	"github.com/canary-debug/kube-vue-admin/tokens"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// LoginOut 退出登录接口
func LoginOut(c *gin.Context) {
	// 使用 Get 而不是 MustGet，防止 panic
	val, exists := c.Get("Authorization")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未找到认证信息，请重新登录",
		})
		return
	}

	// 类型断言检查
	authToken, ok := val.(*jwt.Token)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "Token 格式错误",
		})
		return
	}

	// 从上下文中获取认证token
	Auth := c.MustGet("Authorization").(*jwt.Token)

	// 退出时删除 token
	if authToken.Raw != "" {
		database.RC.Del("valid_token:" + Auth.Raw)
	}
	//database.RC.Del("valid_token:" + Auth.Raw)

	// 解析 token 返回退出用户

	Claims, err := tokens.ReturnClaims(tokens.TokenString)
	if err != nil {
		log.Println(err)
	}
	
	log.Println(Claims)

	// 返回退出成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "退出登录成功",
		"data": gin.H{
			"username": "username",
		},
	})
}
