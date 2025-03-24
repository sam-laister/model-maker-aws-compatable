package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

type JSONMap map[string]interface{}

func (j *JSONMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, j)
}

func (j JSONMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}

type Task struct {
	Id          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Title       string
	Description string
	Completed   bool
	Status      TaskStatus `gorm:"type:TaskStatus"`
	UserId      uint
	Images      []AppFile `gorm:"foreignKey:TaskId"`
	Mesh        *AppFile  `gorm:"foreignKey:TaskId"`
	Metadata    JSONMap   `gorm:"type:json" json:"Metadata"`
}
