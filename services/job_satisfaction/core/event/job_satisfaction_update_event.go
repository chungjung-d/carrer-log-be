package event

import (
	"career-log-be/errors"
	job_satisfaction "career-log-be/models/job_satisfaction"
	"career-log-be/models/job_satisfaction/enums"
	"log"
	"time"

	"gorm.io/gorm"
)

// ConstrainRange는 값이 0-100 범위를 벗어나지 않도록 제한합니다.
func ConstrainRange(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 100 {
		return 100
	}
	return value
}

// ProcessSatisfactionUpdate는 만족도 변경 이벤트를 처리하고 현재 만족도를 업데이트합니다.
func ProcessSatisfactionUpdate(db *gorm.DB, event *job_satisfaction.JobSatisfactionUpdateEvent) error {
	tx := db.Begin()
	if tx.Error != nil {
		return errors.NewInternalError(errors.ErrorCodeDatabaseError, "트랜잭션을 시작할 수 없습니다", tx.Error)
	}

	// 트랜잭션 자동 관리
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 사용자 만족도 객체 준비
	var satisfaction job_satisfaction.UserJobSatisfaction

	// 이벤트 타입에 따른 처리
	if event.EventType == enums.InitEvent {
		// 초기화 이벤트인 경우 새 객체 생성
		satisfaction = job_satisfaction.UserJobSatisfaction{
			UserID: event.UserID,
		}
	} else {
		// 업데이트 이벤트인 경우 기존 데이터 조회
		result := tx.Where("user_id = ?", event.UserID).First(&satisfaction)
		if result.Error != nil {
			tx.Rollback()
			if result.Error == gorm.ErrRecordNotFound {
				return errors.NewNotFoundError(errors.ErrorCodeResourceNotFound, "사용자 만족도가 초기화되지 않았습니다")
			}
			return errors.NewInternalError(errors.ErrorCodeDatabaseError, "사용자 만족도 조회 중 오류가 발생했습니다", result.Error)
		}
	}

	// 만족도 값 업데이트
	satisfaction.Workload = ConstrainRange(satisfaction.Workload + event.Workload)
	satisfaction.Compensation = ConstrainRange(satisfaction.Compensation + event.Compensation)
	satisfaction.Growth = ConstrainRange(satisfaction.Growth + event.Growth)
	satisfaction.WorkEnvironment = ConstrainRange(satisfaction.WorkEnvironment + event.WorkEnvironment)
	satisfaction.WorkRelationships = ConstrainRange(satisfaction.WorkRelationships + event.WorkRelationships)
	satisfaction.WorkValues = ConstrainRange(satisfaction.WorkValues + event.WorkValues)

	// 중요도 값 조회 및 설정
	var importance job_satisfaction.UserJobSatisfactionImportance
	if err := tx.Where("user_id = ?", event.UserID).First(&importance).Error; err == nil {
		satisfaction.WorkloadImportance = importance.Workload
		satisfaction.CompensationImportance = importance.Compensation
		satisfaction.GrowthImportance = importance.Growth
		satisfaction.WorkEnvironmentImportance = importance.WorkEnvironment
		satisfaction.WorkRelationshipsImportance = importance.WorkRelationships
		satisfaction.WorkValuesImportance = importance.WorkValues
	}

	satisfaction.UpdatedAt = time.Now()

	// 데이터베이스 저장 또는 업데이트
	var err error
	if event.EventType == enums.InitEvent {
		err = tx.Create(&satisfaction).Error
	} else {
		err = tx.Save(&satisfaction).Error
	}

	if err != nil {
		tx.Rollback()
		return errors.NewInternalError(errors.ErrorCodeDatabaseError, "사용자 만족도 저장 중 오류가 발생했습니다", err)
	}

	// 이벤트 저장
	if err := tx.Create(event).Error; err != nil {
		tx.Rollback()
		return errors.NewInternalError(errors.ErrorCodeDatabaseError, "이벤트 저장 중 오류가 발생했습니다", err)
	}

	return tx.Commit().Error
}

// PublishJobSatisfactionUpdateEvent는 만족도 업데이트 이벤트를 비동기적으로 처리합니다.
func PublishJobSatisfactionUpdateEvent(db *gorm.DB, event *job_satisfaction.JobSatisfactionUpdateEvent) {
	go func() {
		if err := ProcessSatisfactionUpdate(db, event); err != nil {
			if appErr, ok := err.(*errors.AppError); ok {
				log.Printf("만족도 업데이트 이벤트 처리 실패: %s (Code: %s, Type: %s, Debug: %s)",
					appErr.Message,
					appErr.Code,
					appErr.Type,
					appErr.DebugInfo)
			} else {
				log.Printf("만족도 업데이트 이벤트 처리 중 예상치 못한 오류 발생: %v", err)
			}
		}
	}()
}
