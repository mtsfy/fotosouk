package auth

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/mtsfy/fotosouk/internal/config"
)

func JWTMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:  jwtware.SigningKey{Key: []byte(config.Config("JWT_ACCESS_SECRET"))},
		ContextKey:  "jwt",
		TokenLookup: "cookie:fotosouk_access",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		},
	})
}
