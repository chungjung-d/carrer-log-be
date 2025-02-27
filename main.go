package main

import (
	"career-log-be/config/database"
	"career-log-be/config/env/provider"
	"career-log-be/middleware"
	"career-log-be/models"
	"career-log-be/routes"
	"career-log-be/utils/jwt"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AppContext struct {
	DB       *gorm.DB
	JWT      *jwt.JWTUtils
	Provider provider.EnvProvider
}

func initialize() (*fiber.App, *AppContext, error) {
	// 환경 변수 제공자 초기화
	envProvider := provider.NewDefaultEnvProvider()

	// JWT 유틸리티 초기화
	jwtUtils := jwt.NewJWTUtils(envProvider)

	// 데이터베이스 설정 및 연결
	dbConfig := database.NewConfig()
	db, err := database.NewDatabase(dbConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("could not initialize database connection: %v", err)
	}

	// Auto Migrate
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return nil, nil, fmt.Errorf("could not migrate database: %v", err)
	}

	// Fiber 앱 생성
	app := fiber.New()

	// 컨텍스트 생성
	appCtx := &AppContext{
		DB:       db,
		JWT:      jwtUtils,
		Provider: envProvider,
	}

	// 미들웨어 설정
	app.Use(middleware.DatabaseMiddleware(db))
	app.Use(middleware.JWTMiddleware(jwtUtils))

	// 라우터 설정
	routes.SetupRoutes(app)

	return app, appCtx, nil
}

func main() {
	app, appCtx, err := initialize()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// 기본 DB 인스턴스 가져오기
	sqlDB, err := appCtx.DB.DB()
	if err != nil {
		log.Fatalf("Could not get database instance: %v", err)
	}
	defer sqlDB.Close()

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
