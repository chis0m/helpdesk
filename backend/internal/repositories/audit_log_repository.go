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

func (r *AuditLogRepository) List(ctx context.Context, page, limit int) ([]models.AuditLog, int64, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(&models.AuditLog{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	var rows []models.AuditLog
	err := r.db.WithContext(ctx).Model(&models.AuditLog{}).
		Order("created_at DESC, id DESC").
		Limit(limit).
		Offset(offset).
		Find(&rows).Error
	if err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}
