package routes

import (
	v1 "career-log-be/routes/v1"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// API group
	api := app.Group("/api")

	// Setup v1 routes
	v1.SetupRoutes(api)
}
