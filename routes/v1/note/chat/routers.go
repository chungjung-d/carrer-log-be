package chat

import (
	"career-log-be/middleware"
	chat "career-log-be/services/note/chat"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {
	chatRouter := router.Group("/chat")
	protected := chatRouter.Use(middleware.AuthMiddleware())

	// Get chat by ID
	protected.Get("/:id", chat.HandleGetChat)

	// Stream chat messages
	protected.Post("/:id/stream", chat.HandleStreamChat)
}
