package middleware

import (
	appErrors "career-log-be/errors"
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
			return appErrors.NewAuthorizationError(
				appErrors.ErrorCodeTokenRequired,
				"Authorization header is required",
			)
		}

		// Check if the header has the Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return appErrors.NewAuthorizationError(
				appErrors.ErrorCodeInvalidToken,
				"Invalid authorization header format",
			)
		}

		tokenString := parts[1]

		// Validate token and get claims
		claims, err := jwtUtils.ValidateToken(tokenString)
		if err != nil {
			return appErrors.NewAuthorizationError(
				appErrors.ErrorCodeInvalidToken,
				"Invalid or expired token",
			)
		}

		// Set user information in context
		c.Locals("userID", claims.ID)
		c.Locals("userEmail", claims.Email)

		return c.Next()
	}
}
