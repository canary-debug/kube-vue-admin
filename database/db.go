package database

import (
	"github.com/canary-debug/kube-vue-admin/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "root:1qazWSX!@tcp(43.138.164.205:30001)/kube?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 自动迁移所有模型
	err = db.AutoMigrate(
		&models.Users{},
		// &models.Product{}, // 其他模型可继续添加
	)
	if err != nil {
		panic("failed to migrate database")
	}

	DB = db
}
