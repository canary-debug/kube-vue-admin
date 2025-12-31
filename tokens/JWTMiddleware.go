package tokens

import (
	"net/http"

	"github.com/canary-debug/kube-vue-admin/database"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTMiddleware Gin 中间件：校验 Authorization: Bearer <token>
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 token
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "missing token"})
			return
		}
		tokenStr := header[len("Bearer "):]

		// 解析 token
		var claims Claims
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		// 验证 token 是否过期
		if err != nil || !token.Valid {
			// 删除Redis中过期的token
			database.RC.Del("valid_token:" + tokenStr)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "invalid token"})
			return
		}

		// 额外检查Redis中是否存在该token
		_, err = database.RC.Get("valid_token:" + tokenStr).Result()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "token expired or invalid"})
			return
		}

		// 把 userID 塞进上下文，后续接口直接取
		c.Set("userID", claims.UserID)
		c.Set("Authorization", token)
		c.Next()
	}
}
