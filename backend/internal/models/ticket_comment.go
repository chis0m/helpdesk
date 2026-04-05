package models

import (
	"time"

	"gorm.io/gorm"
)

type TicketComment struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement"`
	TicketID     uint64 `gorm:"not null;index"`
	AuthorUserID uint64 `gorm:"not null;index"`
	Body         string `gorm:"type:text;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type TicketCommentWithAuthor struct {
	ID              uint64    `json:"comment_id"`
	TicketID        uint64    `json:"ticket_id"`
	AuthorUserID    uint64    `json:"author_user_id"`
	AuthorEmail     string    `json:"author_email"`
	AuthorFirstName string    `json:"author_first_name"`
	AuthorLastName  string    `json:"author_last_name"`
	Body            string    `json:"body"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
