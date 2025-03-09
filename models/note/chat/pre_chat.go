package chat

import (
	"career-log-be/utils"
	"time"

	"gorm.io/gorm"
)

const (
	PreChatPrefix = "PRE_CHAT"
)

type PreChat struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (pc *PreChat) BeforeCreate(tx *gorm.DB) error {
	pc.ID = utils.GenerateID(PreChatPrefix)
	return nil
}
