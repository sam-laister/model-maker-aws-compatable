package main

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"

	_ "github.com/joho/godotenv/autoload"

	model "github.com/Soup666/diss-api/model"
)

// Dynamic SQL
type Querier interface {
	// gen.DO // Embeds all default generated methods

	// // SELECT * FROM @@table WHERE name = @name{{if role !=""}} AND role = @role{{end}}
	// FilterWithNameAndRole(name, role string) ([]gen.T, error)

	// // SELECT * FROM @@table WHERE id = @id LIMIT 1
	// GetByID(id int) (*gen.T, error)

	// // SELECT * FROM @@table WHERE email = @email LIMIT 1
	// GetByEmail(email string) (*gen.T, error)

	// // INSERT INTO @@table (email, password) VALUES (@email, @password)
	// InsertValue(email string, password string) (gen.RowsAffected, error) // returns affected rows count and error
}

func main() {

	// log
	log.Println("Connecting to database...")
	log.Println(os.Getenv("DATABASE_URL"))

	g := gen.NewGenerator(gen.Config{
		OutPath: "./query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	gormdb, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	g.UseDB(gormdb)

	// Generate basic type-safe DAO API for struct `model.User` following conventions
	g.ApplyBasic(model.User{})

	// Generate Type Safe API with Dynamic SQL defined on Querier interface for `model.User` and `model.Company`
	// g.ApplyInterface(func(Querier) {}, model.User{})

	// Generate the code
	g.Execute()
}
