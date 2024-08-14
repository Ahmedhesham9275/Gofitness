package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username" binding:"required" gorm:"uniqueIndex;not null"`
	Password string `json:"password" binding:"required" gorm:"not null"`
}
