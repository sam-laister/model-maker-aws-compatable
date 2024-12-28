package model

import "time"

type Mesh struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	Filename string `gorm:"not null"`
	Url      string `gorm:"not null"`

	TaskID uint `gorm:"not null"` // Foreign key
}
