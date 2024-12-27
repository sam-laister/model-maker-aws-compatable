package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectToDatabase establishes a connection to the database and returns a reference to the database instance
func ConnectToDatabase(connectionString string) (*gorm.DB, error) {

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
