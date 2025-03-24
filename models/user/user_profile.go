package user

import (
	"time"
)

type UserProfile struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name" gorm:"not null"`
	Nickname     string    `json:"nickname"`
	Organization string    `json:"organization" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
