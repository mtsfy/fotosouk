package image

import "github.com/gofiber/fiber/v2"

type ImageRepository interface {
	UploadImage(ctx *fiber.Ctx) error
}
