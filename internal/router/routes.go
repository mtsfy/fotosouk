package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mtsfy/fotosouk/internal/auth"
)

func SetupRoutes(app *fiber.App) {
	app.Post("/register", auth.HandleRegister)
	app.Post("/login", auth.HandleLogin)
}
