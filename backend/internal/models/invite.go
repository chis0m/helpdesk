package models

import "time"

type Invite struct {
	ID              uint64     `gorm:"primaryKey;autoIncrement"`
	Email           string     `gorm:"size:120;not null;index"`
	TokenHash       string     `gorm:"column:token_hash;size:64;uniqueIndex;not null"`
	FirstName       string     `gorm:"size:100;not null"`
	LastName        string     `gorm:"size:100;not null"`
	MiddleName      *string    `gorm:"size:100"`
	InvitedByUserID uint64     `gorm:"not null"`
	TargetRole      UserRole   `gorm:"size:20;not null;default:'staff'"`
	ExpiresAt       time.Time  `gorm:"type:datetime(3);not null"`
	UsedAt          *time.Time `gorm:"type:datetime(3)"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (Invite) TableName() string {
	return "invites"
}
