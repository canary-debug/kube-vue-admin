package tokens

import (
	"crypto/rand"
	"os"
)

// jwtKeyFile 把密钥落盘，重启服务不丢失
const jwtKeyFile = ".jwt_key"

// loadOrGenerateKey 优先读文件，没有则生成 32 字节随机密钥
func loadOrGenerateKey() ([]byte, error) {
	//key := make([]byte, 32)
	//if _, err := rand.Read(key); err != nil {
	//	return nil, err
	//}
	//return key, nil
	// 优先从文件读取密钥
	if key, err := os.ReadFile(jwtKeyFile); err == nil {
		return key, nil
	}

	// 文件不存在则生成新密钥并保存
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	if err := os.WriteFile(jwtKeyFile, key, 0600); err != nil {
		return nil, err
	}
	return key, nil

}
