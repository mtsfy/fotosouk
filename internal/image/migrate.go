package image

import (
	"github.com/mtsfy/fotosouk/internal/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.Image{})
}
