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

	images := app.Group("/images")
	images.Use(auth.JWTMiddleware())
	images.Post("/", image.HandleUploadImage(imgSvc))
	images.Get("/", image.HandleGetAllImages(imgSvc))
	images.Post("/:id/transform", image.HandleTransform(imgSvc))
	images.Get("/:id", image.HandleGetImage(imgSvc))
	images.Delete("/:id", image.HandleDeleteImage(imgSvc))
}
