package image

import "github.com/gofiber/fiber/v2"

func HandleUploadImage(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString("Upload")
}
