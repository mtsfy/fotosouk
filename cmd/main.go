package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mtsfy/fotosouk/internal/router"
)

func main() {
	app := fiber.New(fiber.Config{})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("hello world!")
	})

	router.SetupRoutes(app)
	log.Fatal(app.Listen(":3000"))
}
