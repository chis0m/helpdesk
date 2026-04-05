package models

import (
	"time"

	"github.com/google/uuid"
)

type PasswordReset struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement"`
	UserUUID  uuid.UUID  `gorm:"type:char(36);index;not null"`
	TokenHash string     `gorm:"type:char(64);uniqueIndex;not null"`
	ExpiresAt time.Time  `gorm:"type:datetime(3);not null"`
	UsedAt    *time.Time `gorm:"type:datetime(3);index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
