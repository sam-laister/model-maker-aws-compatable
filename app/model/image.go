package model

import "time"

type Image struct {
	ID        uint      `gorm:"primaryKey"`
	Filename  string    `gorm:"not null"`
	Url       string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	TaskID    uint      `gorm:"not null"` // Foreign key
}
