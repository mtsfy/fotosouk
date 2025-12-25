package database

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/mtsfy/fotosouk/internal/config"
	"github.com/mtsfy/fotosouk/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	if url := os.Getenv("DATABASE_URL"); url != "" {
		db, err := gorm.Open(postgres.Open(url), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			log.Fatalf("failed to connect database via DATABASE_URL: %v", err)
		}
		DB = db
	} else {
		p := config.Config("DB_PORT")
		if p == "" {
			p = "5432"
		}
		port, err := strconv.ParseUint(p, 10, 32)
		if err != nil {
			log.Fatalf("failed to parse database port %q: %v", p, err)
		}

		host := config.Config("DB_HOST")
		if host == "" {
			host = "db"
		}
		user := config.Config("DB_USER")
		if user == "" {
			user = "postgres"
		}
		pass := config.Config("DB_PASSWORD")
		name := config.Config("DB_NAME")
		if name == "" {
			name = "postgres"
		}

		dsn := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, pass, name,
		)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			log.Fatalf("failed to connect database: %v", err)
		}
		DB = db
	}

	fmt.Println("Database successfully connected!")

	// retry for migrations
	for i := 0; i < 3; i++ {
		if err := DB.AutoMigrate(&models.User{}); err != nil {
			if i == 2 {
				log.Fatalf("failed to migrate database schemas: %v", err)
			}
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}

	fmt.Println("Database migrated!")
}
