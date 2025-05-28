package database

import (
	"testing"
)

// func TestConnectDatabase(t *testing.T) {

// 	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable TimeZone=%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_TIMEZONE"))
// 	err := ConnectDatabase(dsn)

// 	if err != nil {
// 		t.Fatalf("Failed to connect to database: %v", err)
// 	}

// 	if DB == nil {
// 		t.Fatal("Database connection is nil")
// 	}
// }

func TestAutoMigrate(t *testing.T) {
	err := SetupTestDB(t)

	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
		return
	}

	if DB == nil {
		t.Fatal("Database connection is nil")
	}

	ResetTestDB()
}
