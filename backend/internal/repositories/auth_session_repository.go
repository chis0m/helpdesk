package repositories

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"helpdesk/backend/internal/models"
)

type AuthSessionRepository struct {
	db *gorm.DB
}

func NewAuthSessionRepository(db *gorm.DB) *AuthSessionRepository {
	return &AuthSessionRepository{db: db}
}

func (r *AuthSessionRepository) Create(
	userUUID uuid.UUID,
	sessionID uuid.UUID,
	refreshJTI string,
	refreshExpiresAt time.Time,
) (*models.AuthSession, error) {
	session := &models.AuthSession{
		SessionID:        sessionID,
		UserUUID:         userUUID,
		RefreshJTI:       refreshJTI,
		RefreshExpiresAt: refreshExpiresAt.UTC(),
	}

	if err := r.db.Create(session).Error; err != nil {
		return nil, err
	}

	return session, nil
}

func (r *AuthSessionRepository) GetActiveBySessionID(sessionID uuid.UUID) (*models.AuthSession, error) {
	var session models.AuthSession
	if err := r.db.
		Where("session_id = ? AND revoked_at IS NULL", sessionID).
		First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *AuthSessionRepository) RotateRefreshJTI(
	sessionID uuid.UUID,
	newRefreshJTI string,
	newRefreshExpiresAt time.Time,
) error {
	return r.db.Model(&models.AuthSession{}).
		Where("session_id = ? AND revoked_at IS NULL", sessionID).
		Updates(map[string]any{
			"refresh_jti":        newRefreshJTI,
			"refresh_expires_at": newRefreshExpiresAt.UTC(),
		}).Error
}

func (r *AuthSessionRepository) RevokeBySessionID(sessionID uuid.UUID) error {
	now := time.Now().UTC()
	return r.db.Model(&models.AuthSession{}).
		Where("session_id = ? AND revoked_at IS NULL", sessionID).
		Update("revoked_at", now).Error
}

func (r *AuthSessionRepository) UpsertCSRFToken(
	sessionID uuid.UUID,
	token string,
	expiresAt time.Time,
) error {
	return r.db.Model(&models.AuthSession{}).
		Where("session_id = ? AND revoked_at IS NULL", sessionID).
		Updates(map[string]any{
			"csrf_token":      token,
			"csrf_expires_at": expiresAt.UTC(),
		}).Error
}
