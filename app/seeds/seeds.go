package seeds

import (
	"github.com/Soup666/diss-api/seed"
	"gorm.io/gorm"
)

func All() []seed.Seed {
	return []seed.Seed{
		{
			Name: "CreateTestUser",
			Run: func(db *gorm.DB) error {
				return CreateUser(db, "Test User", "KQmrXe88TwebIMh6AkbEV251Aec2")
			},
		},
		{
			Name: "CreateTestTask",
			Run: func(db *gorm.DB) error {
				return CreateTask(db, "Test task", "Test Description", true, 1)
			},
		},
		{
			Name: "CreateTestFiles",
			Run: func(db *gorm.DB) error {
				return CreateDummyFiles(db)
			},
		},
	}
}
