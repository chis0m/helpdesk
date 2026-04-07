package models

import (
	"time"

	"gorm.io/gorm"
)

type TicketStatus string

const (
	TicketStatusOpen       TicketStatus = "open"
	TicketStatusInProgress TicketStatus = "in_progress"
	TicketStatusResolved   TicketStatus = "resolved"
	TicketStatusClosed     TicketStatus = "closed"
)

type Ticket struct {
	ID             uint64       `gorm:"primaryKey;autoIncrement"`
	ReporterUserID uint64       `gorm:"not null;index"`
	AssignedUserID *uint64      `gorm:"index"`
	Title          string       `gorm:"size:180;not null"`
	Description    string       `gorm:"type:text;not null"`
	Category       string       `gorm:"size:80;not null;default:'general';index"`
	Status         TicketStatus `gorm:"size:30;not null;default:'open';index"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`

	Reporter *User `gorm:"foreignKey:ReporterUserID;references:ID"`
	Assignee *User `gorm:"foreignKey:AssignedUserID;references:ID"`
}
