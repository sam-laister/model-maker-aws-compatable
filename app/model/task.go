package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

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

var TASK_JSON string = `{"Archived":false, "ChatMessages":interface {}(nil), "Completed":false, "CreatedAt":"0001-01-01T00:00:00Z", "DeletedAt":interface {}(nil), "Description":"", "ID":0, "Images":interface {}(nil), "Logs":interface {}(nil), "Mesh":interface {}(nil), "Metadata":interface {}(nil), "Status":"", "Title":"", "UpdatedAt":"0001-01-01T00:00:00Z", "UserId":0}`

type Task struct {
	gorm.Model
	Title        string
	Description  string
	Completed    bool
	Status       TaskStatus `gorm:"type:TaskStatus"`
	UserId       uint
	Images       []AppFile     `gorm:"foreignKey:TaskId"`
	Mesh         *AppFile      `gorm:"foreignKey:TaskId"`
	Metadata     JSONMap       `gorm:"type:json;default:'{}'" json:"Metadata"`
	ChatMessages []ChatMessage `gorm:"foreignKey:TaskId"`
	Logs         []TaskLog     `gorm:"foreignKey:TaskId"`
	Archived     bool          `gorm:"default:false"`
}
