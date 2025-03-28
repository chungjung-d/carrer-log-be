package middleware

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func DatabaseMiddleware(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("db", db)
		return c.Next()
	}
}
