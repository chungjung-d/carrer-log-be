package middleware

import (
	"career-log-be/utils/chatgpt"

	"github.com/gofiber/fiber/v2"
)

func ChatGPTMiddleware(chatGPTService *chatgpt.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("chatgpt", chatGPTService)
		return c.Next()
	}
}
