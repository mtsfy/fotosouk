package auth

import (
	"github.com/gofiber/fiber/v2"
)

type RegisterRequest struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type RegisterResponse struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Username  string `json:"username"`
}

var authService = NewAuthService(&PgRepo{})

func HandleRegister(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	u, err := authService.Register(c.Context(), req.FirstName, req.LastName, req.Email, req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(RegisterResponse{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Username:  u.Username,
	})
}
