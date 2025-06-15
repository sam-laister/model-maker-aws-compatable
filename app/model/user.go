package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email            string
	FirebaseUid      string `gorm:"uniqueIndex"`
	SubscriptionTier string `gorm:"default:free"`
}
