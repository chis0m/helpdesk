package repositories

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"helpdesk/backend/internal/models"
)

type PasswordResetRepository struct {
	db *gorm.DB
}

func NewPasswordResetRepository(db *gorm.DB) *PasswordResetRepository {
	return &PasswordResetRepository{db: db}
}

func (r *PasswordResetRepository) Create(row *models.PasswordReset) error {
	return r.db.Create(row).Error
}

func (r *PasswordResetRepository) InvalidateUnusedForUser(userUUID uuid.UUID, at time.Time) error {
	return r.db.Model(&models.PasswordReset{}).
		Where("user_uuid = ? AND used_at IS NULL", userUUID).
		Update("used_at", at.UTC()).Error
}

func (r *PasswordResetRepository) GetByTokenHash(tokenHash string) (*models.PasswordReset, error) {
	var row models.PasswordReset
	if err := r.db.Where("token_hash = ?", tokenHash).First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *PasswordResetRepository) MarkUsed(id uint64, at time.Time) error {
	return r.db.Model(&models.PasswordReset{}).
		Where("id = ?", id).
		Update("used_at", at.UTC()).Error
}
