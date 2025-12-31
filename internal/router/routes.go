package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mtsfy/fotosouk/internal/auth"
	"github.com/mtsfy/fotosouk/internal/image"
)

func SetupRoutes(app *fiber.App, imgSvc *image.ImageService, authSvc *auth.AuthService) {
	app.Post("/register", auth.HandleRegister(authSvc))
	app.Post("/login", auth.HandleLogin(authSvc))
	app.Post("/refresh", auth.HandleRefresh(authSvc))

	protected := app.Group("/images")
	protected.Use(auth.JWTMiddleware())
	protected.Post("/", image.HandleUploadImage(imgSvc))
	protected.Get("/", image.HandleGetAllImages(imgSvc))
	protected.Post("/:id/transform", image.HandleTransform(imgSvc))
	protected.Get("/:id", image.HandleGetImage(imgSvc))
}
