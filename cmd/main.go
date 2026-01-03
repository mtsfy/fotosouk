package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mtsfy/fotosouk/internal/auth"
	"github.com/mtsfy/fotosouk/internal/config"
	"github.com/mtsfy/fotosouk/internal/database"
	"github.com/mtsfy/fotosouk/internal/image"
	"github.com/mtsfy/fotosouk/internal/router"
	"github.com/mtsfy/fotosouk/internal/storage"
	"github.com/mtsfy/fotosouk/internal/transformer"
)

func main() {
	app := fiber.New(fiber.Config{})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("hello world!")
	})

	database.Connect()
	if err := auth.Migrate(database.DB); err != nil {
		log.Fatalf("failed to migrate auth: %v", err)
	}
	if err := image.Migrate(database.DB); err != nil {
		log.Fatalf("failed to migrate image: %v", err)
	}

	stor, err := storage.NewS3Storage(
		config.Config("AWS_S3_BUCKET"),
		config.Config("AWS_REGION"),
	)
	if err != nil {
		log.Fatalf("failed to init storage: %v", err)
	}

	imgSvc := image.NewImageService(&image.PgRepo{}, stor, &transformer.GoImageTransformer{})
	authSvc := auth.NewAuthService(&auth.PgRepo{})

	router.SetupRoutes(app, imgSvc, authSvc)

	port := config.Config("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Starting server on :%s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
