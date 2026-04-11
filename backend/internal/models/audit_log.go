package models

import (
	"time"
)

// AuditLog is append-only.
type AuditLog struct {
	ID           uint64    `gorm:"column:id;primaryKey"`
	ActorUserID  *uint64   `gorm:"column:actor_user_id"`
	SessionID    *string   `gorm:"column:session_id;size:36"`
	TokenJTI     *string   `gorm:"column:token_jti;size:64"`
	HTTPMethod   string    `gorm:"column:http_method;size:16"`
	Path         string    `gorm:"column:path;size:512"`
	Action       string    `gorm:"column:action;size:160"`
	Success      bool      `gorm:"column:success"`
	ErrorCode    *string   `gorm:"column:error_code;size:96"`
	ResourceType *string   `gorm:"column:resource_type;size:64"`
	ResourceID   *uint64   `gorm:"column:resource_id"`
	IP           *string   `gorm:"column:ip;size:64"`
	UserAgent    *string   `gorm:"column:user_agent;size:512"`
	Metadata     []byte    `gorm:"column:metadata;type:json"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
