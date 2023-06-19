package models

import (
	"github.com/google/uuid"
	"time"
)

type Token struct {
	ID           uuid.UUID `json:"id" gorm:"primarykey"`
	UserID       int       `json:"user_id"`
	RefreshToken string    `gorm:"size:2048" json:"refresh_token"`
	PlatformID   int       `gorm:"not null" json:"platform_id"`
	IsBlocked    bool      `gorm:"not null;default:false" json:"is_blocked"`
	ExpiresAt    time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (t *Token) TableName() string {
	return "tokens"
}
