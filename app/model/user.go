package model

import (
	"time"
)

type User struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
	Email       string
	FirebaseUid string `gorm:"uniqueIndex"`
}
