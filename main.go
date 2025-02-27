package main

import (
	"career-log-be/config/database"
	"career-log-be/middleware"
	"career-log-be/models"
	"career-log-be/routes"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// 데이터베이스 설정
	dbConfig := database.NewConfig()

	// 데이터베이스 연결
	db, err := database.NewDatabase(dbConfig)
	if err != nil {
		log.Fatalf("Could not initialize database connection: %v", err)
	}

	// 기본 DB 인스턴스 가져오기
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Could not get database instance: %v", err)
	}
	defer sqlDB.Close()

	// Auto Migrate
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Could not migrate database: %v", err)
	}

	app := fiber.New()

	// 전역 미들웨어로 DB 설정
	app.Use(middleware.DatabaseMiddleware(db))

	// 라우터 설정
	routes.SetupRoutes(app)

	// 그레이스풀 셧다운을 위한 채널 설정
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Gracefully shutting down...")
		_ = app.Shutdown()
	}()

	// 서버 시작
	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
