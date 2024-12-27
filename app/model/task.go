package model

import (
	"time"
)

type Task struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
	Title       string
	Description string
	Completed   bool
	UserID      uint
}
