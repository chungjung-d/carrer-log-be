package user

import (
	"career-log-be/middleware"
	"career-log-be/services/user"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {
	userRouter := router.Group("/user")

	// 보호된 라우트 그룹
	protected := userRouter.Use(middleware.AuthMiddleware())

	// 프로필 생성
	protected.Post("/profile", user.HandleCreateUserProfile())

	// // 프로필 조회
	// protected.Get("/profile", userService.HandleGetProfile())

	// // 프로필 업데이트
	// protected.Put("/profile", userService.HandleUpdateProfile())
}
