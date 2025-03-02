package user

import (
	"career-log-be/utils"
	"time"

	"gorm.io/gorm"
)

const (
	UserPrefix = "USR"
)

type User struct {
	ID        string `gorm:"primaryKey;type:varchar(100)"`
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	user.ID = utils.GenerateID(UserPrefix)
	return nil
}
