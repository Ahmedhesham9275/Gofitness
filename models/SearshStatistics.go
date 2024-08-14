package models

import (
	"sync"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type SearchStatistic struct {
	gorm.Model
	Word        string `gorm:"uniqueIndex;not null"`
	SearchCount uint   `gorm:"default:0"`
	LastTF      uint   `gorm:"default:0"`
	LastDF      uint   `gorm:"default:0"`
	History     datatypes.JSON
}

type TFDFHistory struct {
	TF        int       `json:"tf"`
	DF        int       `json:"df"`
	Timestamp time.Time `json:"timestamp"`
}

// Mutex for SearchStatistic to handle concurrency
var Mutex sync.Mutex
