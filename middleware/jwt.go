package middleware

import (
	"career-log-be/utils/jwt"

	"github.com/gofiber/fiber/v2"
)

func JWTMiddleware(jwtUtils *jwt.JWTUtils) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("jwt", jwtUtils)
		return c.Next()
	}
}
