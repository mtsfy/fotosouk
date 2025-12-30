package image

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func HandleUploadImage(svc *ImageService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		token, ok := c.Locals("jwt").(*jwt.Token)
		if !ok {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		claims := token.Claims.(jwt.MapClaims)

		uid, ok := claims["user_id"].(float64)
		if !ok {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		fh, err := c.FormFile("image")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "file is required",
			})
		}

		userID := int(uid)
		file, err := fh.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to open file",
			})
		}
		defer file.Close()

		filename := fh.Filename
		fileSize := fh.Size
		mimeType := fh.Header.Get("Content-Type")

		img, err := svc.UploadImage(c.Context(), userID, filename, file, fileSize, mimeType)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"id":         img.ID,
			"width":      img.Width,
			"height":     img.Height,
			"filename":   img.Filename,
			"url":        img.Url,
			"mime_type":  img.MimeType,
			"file_size":  img.FileSize,
			"created_at": img.CreatedAt,
		})
	}
}

func HandleGetAllImages(svc *ImageService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.Locals("jwt").(*jwt.Token)
		if !ok {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		uid, ok := claims["user_id"].(float64)
		if !ok {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		userID := int(uid)
		images, err := svc.GetAllImages(c.Context(), userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"images": images,
		})
	}
}

func HandleGetImage(svc *ImageService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.Locals("jwt").(*jwt.Token)
		if !ok {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		uid, ok := claims["user_id"].(float64)
		if !ok {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		userID := int(uid)
		id := c.Params("id")
		imgID, err := strconv.Atoi(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid image id",
			})
		}

		img, err := svc.GetImageDetail(c.Context(), userID, imgID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "image not found",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"image": img,
		})
	}
}
