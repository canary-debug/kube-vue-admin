package tokens

import (
	"encoding/base64"
	"log"
)

func LoadToekn() {
	// 加载密钥
	var err error
	jwtKey, err = loadOrGenerateKey()
	if err != nil {
		log.Fatalf("load key: %v", err)
	}
	base64.StdEncoding.EncodeToString(jwtKey)

}
