// Auth routes

package auth

import (
	authService "career-log-be/services/auth"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {
	auth := router.Group("/auth")

	// Register route
	auth.Post("/register", authService.HandleRegister())
}
