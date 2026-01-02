package database

import (
	"fmt"
	"log"
	"os"

	"github.com/canary-debug/kube-vue-admin/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// 从环境变量读取 MySQL user password
	MySQL_User := os.Getenv("MYSQL_USER")
	MySQL_Password := os.Getenv("MYSQL_PASSWORD")
	MySQL_Host := os.Getenv("MYSQL_HOST")
	MySQL_Port := os.Getenv("MYSQL_PORT")
	MySQL_DBName := os.Getenv("MYSQL_DBNAME")

	// 动态拼接 DSN
	// 格式: 用户名:密码@tcp(地址)/数据库名?参数
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		MySQL_User,
		MySQL_Password,
		MySQL_Host,
		MySQL_Port,
		MySQL_DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("failed to connect database: %v", err)
		//panic("failed to connect database")
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
