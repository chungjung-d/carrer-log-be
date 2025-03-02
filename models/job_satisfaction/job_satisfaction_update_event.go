package job_satisfaction

import (
	"career-log-be/models/job_satisfaction/enums"
	"career-log-be/utils"

	"time"

	"gorm.io/gorm"
)

const (
	// JobSatisfactionEventPrefix는 직무 만족도 변경 이벤트 ID의 접두사입니다
	JobSatisfactionEventPrefix = "JOB_SAT_UPDATE_EVENT"
)

type JobSatisfactionUpdateEvent struct {
	ID                string                               `json:"id" gorm:"primaryKey;not null"`
	UserID            string                               `json:"userId" gorm:"index;not null"`
	EventType         enums.JobSatisfactionUpdateEventType `json:"eventType" gorm:"type:varchar(20);not null"`
	Workload          int                                  `json:"workload" gorm:"not null"` // 변화량 (양수 또는 음수)
	Compensation      int                                  `json:"compensation" gorm:"not null"`
	Growth            int                                  `json:"growth" gorm:"not null"`
	WorkEnvironment   int                                  `json:"workEnvironment" gorm:"column:work_environment;not null"`
	WorkRelationships int                                  `json:"workRelationships" gorm:"column:work_relationships;not null"`
	WorkValues        int                                  `json:"workValues" gorm:"column:work_values;not null"`
	SourceId          *string                              `json:"sourceId"` // 참조 ID (S3에 저장된 대화 내용 참조), nullable
	CreatedAt         time.Time                            `json:"createdAt" gorm:"not null"`
}

func (u *JobSatisfactionUpdateEvent) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = utils.GenerateID(JobSatisfactionEventPrefix)
	return
}
