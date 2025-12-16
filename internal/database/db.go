package database

import (
	"fmt"
	"strconv"

	"github.com/mtsfy/fotosouk/internal/config"
	"github.com/mtsfy/fotosouk/internal/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	p := config.Config("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		panic("failed to parse database port")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Config("DB_HOST"),
		port,
		config.Config("DB_USER"),
		config.Config("DB_PASSWORD"),
		config.Config("DB_NAME"),
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Database successfully connected!")

	err = DB.AutoMigrate(&user.User{})

	if err != nil {
		panic("failed to migrate database schemas")
	}

	fmt.Println("Database migrated!")
}
