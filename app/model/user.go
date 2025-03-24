package model

import (
	"time"
)

type User struct {
	Id          uint      `gorm:"primaryKey"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	DeletedAt   *time.Time
	Email       string
	FirebaseUid string `gorm:"uniqueIndex"`
}
