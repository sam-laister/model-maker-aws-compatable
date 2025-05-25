package model

import (
	"time"

	"gorm.io/gorm"
)

type TaskLog struct {
	gorm.Model
	CreatedAt time.Time `gorm:"autoCreateTime"`
	Message   string    `gorm:"type:text;not null"`
	TaskId    uint      `gorm:"not null"` // Foreign key
}
