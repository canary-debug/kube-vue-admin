package tokens

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

// 解密函数 (token)
func ParseJwt(key any, jwtStr string, options ...jwt.ParserOption) (jwt.Claims, error) {
	token, err := jwt.Parse(jwtStr, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	}, options...)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token is invalid")
	}
	return token.Claims, nil
}

// 解密流程函数
func ReturnClaims(jwtStr string) (jwt.Claims, error) {
	// 入参：jwtStr 是真实的 JWT 字符串（如从请求头获取的）
	claims, err := ParseJwt(
		jwtKey,                                  // 1. 签名密钥
		jwtStr,                                  // 2. 真实的 JWT 字符串（核心修正）
		jwt.WithValidMethods([]string{"HS256"}), // 3. 解析选项：仅允许 HS256
	)
	if err != nil {
		return nil, err
	}
	return claims, nil
}
