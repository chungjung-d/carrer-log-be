package chat

import (
	"career-log-be/middleware"
	// chat "career-log-be/services/chat"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {

	chatRouter := router.Group("/chat")

	protected := chatRouter.Use(middleware.AuthMiddleware())

	// protected.Post("/", chat.HandleCreateChat())
}
