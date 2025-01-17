package model

import (
	"time"
)

type User struct {
	ID          uint      `gorm:"primaryKey"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	DeletedAt   *time.Time
	Email       string
	FirebaseUid string `gorm:"uniqueIndex"`
}
