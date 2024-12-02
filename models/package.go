package models

import (
	"gorm.io/gorm"
)

type Package struct {
	gorm.Model
	Title       string `json:"title" binding:"required"`
	Description string `json:"Description" binding:"required"`
	UserID      uint
}
