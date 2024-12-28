package model

import (
	"time"
)

type Task struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time  `gorm:"autoCreateTime"`
	DeletedAt   *time.Time `gorm:"autoUpdateTime"`
	Title       string
	Description string
	Completed   bool
	UserID      uint
	Images      []Image `gorm:"foreignKey:TaskID"`
}
