package note

import (
	"career-log-be/routes/v1/note/chat"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {
	noteRouter := router.Group("/note")

	chat.SetupRoutes(noteRouter)

}
