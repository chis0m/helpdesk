package repositories

import (
	"context"

	"gorm.io/gorm"

	"helpdesk/backend/internal/models"
)

type AuditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) Create(ctx context.Context, row *models.AuditLog) error {
	return r.db.WithContext(ctx).Create(row).Error
}
