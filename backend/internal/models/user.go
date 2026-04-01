package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleStaff UserRole = "staff"
	RoleAdmin UserRole = "admin"
)

type User struct {
	ID                 uint64     `gorm:"primaryKey;autoIncrement"`
	UUID               uuid.UUID  `gorm:"type:char(36);uniqueIndex;not null"`
	Email              string     `gorm:"size:120;uniqueIndex;not null"`
	PasswordHash       string     `gorm:"size:255;not null"`
	FirstName          string     `gorm:"size:100;not null"`
	LastName           string     `gorm:"size:100;not null"`
	MiddleName         *string    `gorm:"size:100"`
	Role               UserRole   `gorm:"size:20;not null;default:'user'"`
	IsActive           bool       `gorm:"not null;default:true"`
	MustChangePassword bool       `gorm:"not null;default:true"`
	PasswordChangedAt  *time.Time `gorm:"type:datetime(3)"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt `gorm:"index"`
}
