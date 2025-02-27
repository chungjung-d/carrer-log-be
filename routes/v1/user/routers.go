package user

import (
	"career-log-be/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {
	user := router.Group("/user")

	// 보호된 라우트 그룹
	protected := user.Use(middleware.AuthMiddleware())

	// // 프로필 조회
	// protected.Get("/profile", userService.HandleGetProfile())

	// // 프로필 업데이트
	// protected.Put("/profile", userService.HandleUpdateProfile())
}
