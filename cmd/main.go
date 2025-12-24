package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mtsfy/fotosouk/internal/auth"
	"github.com/mtsfy/fotosouk/internal/config"
	"github.com/mtsfy/fotosouk/internal/database"
	"github.com/mtsfy/fotosouk/internal/router"
)

func main() {
	app := fiber.New(fiber.Config{})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("hello world!")
	})

	database.Connect()
	if err := auth.Migrate(database.DB); err != nil {
		log.Fatalf("failed to migrate auth schemas: %v", err)
	}

	router.SetupRoutes(app)

	port := config.Config("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Starting server on :%s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
