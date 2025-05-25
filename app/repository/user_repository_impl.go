package repositories

import (
	"github.com/Soup666/diss-api/database"
	models "github.com/Soup666/diss-api/model"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{DB: db}
}

func (repo *UserRepositoryImpl) GetUserFromFirebaseUID(apiKey string) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("firebase_uid = ?", apiKey).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepositoryImpl) Create(user *models.User) error {
	return repo.DB.Create(&user).Error
}

func (repo *UserRepositoryImpl) UpdateUser(user *models.User) error {
	return repo.DB.Save(&user).Error
}

func (repo *UserRepositoryImpl) GetUsers() ([]*models.User, error) {
	var users []*models.User
	if err := database.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *UserRepositoryImpl) DeleteUser(user *models.User) error {
	return repo.DB.Delete(&user).Error
}
