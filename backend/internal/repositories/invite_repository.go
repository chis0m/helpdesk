package repositories

import (
	"time"

	"gorm.io/gorm"

	"helpdesk/backend/internal/models"
)

type InviteRepository struct {
	db *gorm.DB
}

func NewInviteRepository(db *gorm.DB) *InviteRepository {
	return &InviteRepository{db: db}
}

func (r *InviteRepository) Create(inv *models.Invite) error {
	return r.db.Create(inv).Error
}

func (r *InviteRepository) GetByTokenHash(tokenHash string) (*models.Invite, error) {
	var inv models.Invite
	if err := r.db.Where("token_hash = ?", tokenHash).First(&inv).Error; err != nil {
		return nil, err
	}
	return &inv, nil
}

func (r *InviteRepository) HasPendingInviteForEmail(email string, now time.Time) (bool, error) {
	var count int64
	err := r.db.Model(&models.Invite{}).
		Where("email = ? AND used_at IS NULL AND expires_at > ?", email, now).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *InviteRepository) MarkUsed(id uint64, usedAt time.Time) error {
	return r.db.Model(&models.Invite{}).
		Where("id = ?", id).
		Update("used_at", usedAt).Error
}
