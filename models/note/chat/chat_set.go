package chat

import (
	"career-log-be/utils"
	"time"

	"gorm.io/gorm"
)

const (
	ChatSetPrefix = "CH_SET"
)

type ChatSet struct {
	ID        string         `gorm:"primaryKey;type:varchar(100)" json:"id"`
	UserID    string         `gorm:"type:varchar(100);not null" json:"user_id"`
	Title     string         `json:"title"`
	ChatData  ChatData       `gorm:"type:jsonb" json:"chat_data"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (chat *ChatSet) BeforeCreate(tx *gorm.DB) error {
	chat.ID = utils.GenerateID(ChatSetPrefix)
	return nil
}
