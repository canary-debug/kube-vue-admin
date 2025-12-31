package tokens

import (
	"time"

	"github.com/canary-debug/kube-vue-admin/database"
	"github.com/golang-jwt/jwt/v5"
)

// jwtKey 密钥
var jwtKey []byte

// 全局可以调用的 Token
var TokenString string

type Claims struct {
	UserID   uint   `json:"id"`
	UserName string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 为指定 userID(ID) 签发 1h 有效期 token
func GenerateToken(userID uint, userName string) (string, error) {
	claims := Claims{
		UserID:   userID,
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)), // 过期时间 1h
			IssuedAt:  jwt.NewNumericDate(time.Now()),                // 当前时间
			Issuer:    "gin-kube",                                    // 签发者
		},
	}

	// 生成 token
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtKey)
	TokenString = tokenString
	if err != nil {
		return "", err
	}

	/*
		将token存储到Redis，设置相同过期时间
		密钥=key
		用户ID=value
	*/
	database.RC.Set("valid_token:"+tokenString, userID, time.Hour)

	return tokenString, nil
}
