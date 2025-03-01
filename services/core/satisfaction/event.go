package satisfaction

import (
	"career-log-be/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

// ConstrainRange는 값이 0-100 범위를 벗어나지 않도록 제한합니다.
func ConstrainRange(value int) int {
	if value < 0 {
		return 0
	}
	if value > 100 {
		return 100
	}
	return value
}

// ProcessSatisfactionUpdate는 만족도 변경 이벤트를 처리하고 현재 만족도를 업데이트합니다.
func ProcessSatisfactionUpdate(db *gorm.DB, event *models.JobSatisfactionUpdateEvent) error {
	// 트랜잭션 시작
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 롤백 함수
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 사용자 만족도 조회
	var satisfaction models.UserJobSatisfaction
	result := tx.Where("user_id = ?", event.UserID).First(&satisfaction)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 만족도 레코드가 없는 경우 에러 반환 (초기화 필요)
			tx.Rollback()
			return errors.New("사용자 만족도가 초기화되지 않았습니다")
		}
		tx.Rollback()
		return result.Error
	}

	// 만족도 업데이트
	satisfaction.Workload = ConstrainRange(satisfaction.Workload + event.Workload)
	satisfaction.Compensation = ConstrainRange(satisfaction.Compensation + event.Compensation)
	satisfaction.Growth = ConstrainRange(satisfaction.Growth + event.Growth)
	satisfaction.WorkEnvironment = ConstrainRange(satisfaction.WorkEnvironment + event.WorkEnvironment)
	satisfaction.WorkRelationships = ConstrainRange(satisfaction.WorkRelationships + event.WorkRelationships)
	satisfaction.WorkValues = ConstrainRange(satisfaction.WorkValues + event.WorkValues)
	satisfaction.UpdatedAt = time.Now()

	// 이벤트 저장
	if err := tx.Create(event).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 만족도 업데이트
	if err := tx.Save(&satisfaction).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 트랜잭션 커밋
	return tx.Commit().Error
}

// // GetUserSatisfaction은 사용자의 현재 만족도 정보를 조회합니다.
// func GetUserSatisfaction(db *gorm.DB, userID string) (*models.UserJobSatisfaction, error) {
// 	var satisfaction models.UserJobSatisfaction
// 	result := db.Where("user_id = ?", userID).First(&satisfaction)

// 	if result.Error != nil {
// 		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 			return nil, errors.New("사용자 만족도 정보가 없습니다")
// 		}
// 		return nil, result.Error
// 	}

// 	return &satisfaction, nil
// }

// // GetUserSatisfactionHistory는 사용자의 만족도 변경 이력을 조회합니다.
// func GetUserSatisfactionHistory(db *gorm.DB, userID string, limit int) ([]*models.JobSatisfactionUpdateEvent, error) {
// 	var events []*models.JobSatisfactionUpdateEvent

// 	result := db.Where("user_id = ?", userID).
// 		Order("created_at DESC").
// 		Limit(limit).
// 		Find(&events)

// 	if result.Error != nil {
// 		return nil, result.Error
// 	}

// 	return events, nil
// }
