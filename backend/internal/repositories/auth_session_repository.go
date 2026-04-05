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
	userAgent *string,
	ip *string,
) (*models.AuthSession, error) {
	session := &models.AuthSession{
		SessionID:        sessionID,
		UserUUID:         userUUID,
		RefreshJTI:       refreshJTI,
		RefreshExpiresAt: refreshExpiresAt.UTC(),
		UserAgent:        userAgent,
		IP:               ip,
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

func (r *AuthSessionRepository) ListActiveByUserUUID(userUUID uuid.UUID) ([]models.AuthSession, error) {
	var rows []models.AuthSession
	if err := r.db.
		Where("user_uuid = ? AND revoked_at IS NULL", userUUID).
		Order("created_at DESC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// RevokeActiveSessionForUser revokes one session if it belongs to the user. Returns whether a row was updated.
func (r *AuthSessionRepository) RevokeActiveSessionForUser(sessionID uuid.UUID, userUUID uuid.UUID) (bool, error) {
	now := time.Now().UTC()
	res := r.db.Model(&models.AuthSession{}).
		Where("session_id = ? AND user_uuid = ? AND revoked_at IS NULL", sessionID, userUUID).
		Update("revoked_at", now)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

func (r *AuthSessionRepository) RevokeOtherActiveSessionsForUser(currentSessionID uuid.UUID, userUUID uuid.UUID) error {
	now := time.Now().UTC()
	return r.db.Model(&models.AuthSession{}).
		Where("user_uuid = ? AND session_id != ? AND revoked_at IS NULL", userUUID, currentSessionID).
		Update("revoked_at", now).Error
}

func (r *AuthSessionRepository) RevokeAllActiveForUser(userUUID uuid.UUID) error {
	now := time.Now().UTC()
	return r.db.Model(&models.AuthSession{}).
		Where("user_uuid = ? AND revoked_at IS NULL", userUUID).
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
