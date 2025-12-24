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

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Token    string `json:"token"`
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

func HandleLogin(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	u, token, err := authService.Login(c.Context(), req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "fotosouk_access",
		Value:    token.AccessToken,
		MaxAge:   24 * 60 * 60,
		Path:     "/",
		Secure:   true,
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteStrictMode,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "fotosouk_refresh",
		Value:    token.RefreshToken,
		MaxAge:   24 * 60 * 60 * 7,
		Path:     "/",
		Secure:   true,
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteStrictMode,
	})

	return c.Status(fiber.StatusCreated).JSON(LoginResponse{
		ID:       u.ID,
		Email:    u.Email,
		Username: u.Username,
		Token:    token.AccessToken,
	})
}

func HandleRefresh(c *fiber.Ctx) error {
	refreshToken := c.Cookies("fotosouk_refresh")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "refresh token required",
		})
	}

	token, err := authService.RefreshAccessToken(c.Context(), refreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "fotosouk_access",
		Value:    token.AccessToken,
		MaxAge:   24 * 60 * 60,
		Path:     "/",
		Secure:   true,
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteStrictMode,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"accessToken": token.AccessToken,
	})
}
