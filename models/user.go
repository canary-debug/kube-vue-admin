package models

import "gorm.io/gorm"

// User 用户模型
type Users struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
