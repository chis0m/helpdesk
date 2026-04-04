package models

import (
	"time"

	"github.com/google/uuid"
)

type AuthSession struct {
	ID               uint64     `gorm:"primaryKey;autoIncrement"`
	SessionID        uuid.UUID  `gorm:"type:char(36);uniqueIndex;not null"`
	UserUUID         uuid.UUID  `gorm:"type:char(36);index;not null"`
	RefreshJTI       string     `gorm:"type:char(36);index;not null"`
	RefreshExpiresAt time.Time  `gorm:"type:datetime(3);not null"`
	CSRFToken        *string    `gorm:"type:varchar(128)"`
	CSRFExpiresAt    *time.Time `gorm:"type:datetime(3)"`
	RevokedAt        *time.Time `gorm:"type:datetime(3);index"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
