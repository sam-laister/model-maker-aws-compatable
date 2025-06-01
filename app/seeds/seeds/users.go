package seeds

import (
	"github.com/Soup666/modelmaker/model"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, email string, firebaseUid string) error {
	if err := db.Create(&model.User{Email: email, FirebaseUid: firebaseUid}).Error; err != nil {
		return nil
	}
	return nil
}
