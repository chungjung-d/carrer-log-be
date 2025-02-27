package middleware

import (
	"career-log-be/utils/jwt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware checks for valid JWT token and sets user info in context
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get JWT utils from context
		jwtUtils := c.Locals("jwt").(*jwt.JWTUtils)

		// Get token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header is required",
			})
		}

		// Check if the header has the Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		tokenString := parts[1]

		// Validate token and get claims
		claims, err := jwtUtils.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Set user information in context
		c.Locals("userId", claims.ID)
		c.Locals("userEmail", claims.Email)

		return c.Next()
	}
}
