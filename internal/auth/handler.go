package auth

import (
	"github.com/gofiber/fiber/v2"
)

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func HandleRegister(c *fiber.Ctx) error {
	var auth Auth
	if err := c.BodyParser(&auth); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"message": "register"})
}

func HandleLogin(c *fiber.Ctx) error {
	var auth Auth
	if err := c.BodyParser(&auth); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"message": "login"})
}
