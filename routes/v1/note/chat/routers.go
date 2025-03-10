package chat

import (
	"career-log-be/middleware"
	chat "career-log-be/services/note/chat"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {
	chatRouter := router.Group("/chat")
	protected := chatRouter.Use(middleware.AuthMiddleware())

	// Get all pre-chats
	protected.Get("/pre-chats", chat.HandleListPreChats)

	// Create new pre-chat
	protected.Post("/pre-chats", chat.HandleCreatePreChat)

	// Create new chat
	protected.Post("/create", chat.HandleCreateChat)

	// Analyze daily chat manually (Scheduler Test API)
	protected.Post("/analyze-daily", chat.HandleAnalyzeDailyChat)

	// Get chat by ID
	protected.Get("/:id", chat.HandleGetChat)

	// Stream chat messages
	protected.Post("/:id", chat.HandleChat)
}
