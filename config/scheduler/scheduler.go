package scheduler

import (
	"career-log-be/services/note/chat/scheduler"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// InitSchedulers 모든 스케줄러를 초기화하는 함수
func InitSchedulers(app *fiber.App, db *gorm.DB) error {
	// 채팅 분석 스케줄러 초기화
	if err := scheduler.InitChatAnalyzeScheduler(app, db); err != nil {
		return err
	}

	return nil
}
