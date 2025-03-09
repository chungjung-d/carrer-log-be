package main

import (
	"career-log-be/config/database"
	"career-log-be/middleware"
	job_satisfaction "career-log-be/models/job_satisfaction"
	user "career-log-be/models/user"
	"career-log-be/routes"
	"career-log-be/utils/jwt"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type AppContext struct {
	DB  *gorm.DB
	JWT *jwt.JWTUtils
}

func initialize() (*fiber.App, *AppContext, error) {
	// .env 파일 로드
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// JWT 유틸리티 초기화
	jwtUtils := jwt.NewJWTUtils()

	// 데이터베이스 설정 및 연결
	dbConfig := database.NewConfig()
	db, err := database.NewDatabase(dbConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("could not initialize database connection: %v", err)
	}

	// Auto Migrate
	if err := db.AutoMigrate(&user.User{}, &user.UserProfile{}, &job_satisfaction.UserJobSatisfactionImportance{}, &job_satisfaction.UserJobSatisfaction{}, &job_satisfaction.JobSatisfactionUpdateEvent{}); err != nil {
		return nil, nil, fmt.Errorf("could not migrate database: %v", err)
	}

	// Fiber 앱 생성 (에러 핸들러 등록)
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler(),
	})

	// 컨텍스트 생성
	appCtx := &AppContext{
		DB:  db,
		JWT: jwtUtils,
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
