package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config 首字母大写，导出结构体
type Config struct {
	Weather struct {
		// 字段首字母大写，导出
		ApiKey     string `yaml:"api_key"`
		DistrictID string `yaml:"district_id"`
	} `yaml:"weather"`
}

// C 首字母大写，导出的包级变量，其他包可直接引用
var C Config

// Init 初始化配置（在项目启动时调用）
func Init() error {
	file, err := os.ReadFile("configs/config.yaml")
	if err != nil {
		return err
	}
	return yaml.Unmarshal(file, &C)
}
