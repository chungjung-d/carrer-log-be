package v1

import (
	"career-log-be/routes/v1/auth"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {
	v1 := router.Group("/v1")

	// Setup auth routes
	auth.SetupRoutes(v1)
}
