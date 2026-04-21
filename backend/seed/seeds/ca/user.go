package ca

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"helpdesk/backend/internal/auth"
	"helpdesk/backend/internal/models"
)

const (
	EmailMustChange  = "must.change@company-a.com"
	EmailMarkAnthony = "mark.anthony@company-a.com"
	EmailJaneDoe     = "jane.doe@company-b.com"
	EmailAlexJones   = "alex.jones@company-b.com"
)

func EnsureCustomerUsers(db *gorm.DB) (*models.User, *models.User, *models.User, *models.User, error) {
	hash, err := auth.HashPassword(caTestPassword)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	settled := time.Now().UTC().Add(-48 * time.Hour)

	uMust, _, err := firstOrCreateUser(
		db,
		EmailMustChange,
		"Must",
		"Change",
		models.RoleUser,
		true,
		nil,
		hash,
	)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if err := ensureArgon2idPasswordHash(db, uMust, caTestPassword); err != nil {
		return nil, nil, nil, nil, err
	}

	uMark, _, err := firstOrCreateUser(
		db,
		EmailMarkAnthony,
		"Mark",
		"Anthony",
		models.RoleUser,
		false,
		&settled,
		hash,
	)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if err := ensureArgon2idPasswordHash(db, uMark, caTestPassword); err != nil {
		return nil, nil, nil, nil, err
	}

	uJane, _, err := firstOrCreateUser(
		db,
		EmailJaneDoe,
		"Jane",
		"Doe",
		models.RoleUser,
		false,
		&settled,
		hash,
	)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if err := ensureArgon2idPasswordHash(db, uJane, caTestPassword); err != nil {
		return nil, nil, nil, nil, err
	}

	uAlex, _, err := firstOrCreateUser(
		db,
		EmailAlexJones,
		"Alex",
		"Jones",
		models.RoleUser,
		false,
		&settled,
		hash,
	)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if err := ensureArgon2idPasswordHash(db, uAlex, caTestPassword); err != nil {
		return nil, nil, nil, nil, err
	}

	if uMust == nil || uMark == nil || uJane == nil || uAlex == nil {
		return nil, nil, nil, nil, fmt.Errorf("EnsureCustomerUsers: missing user(s)")
	}
	return uMust, uMark, uJane, uAlex, nil
}

func ensureArgon2idPasswordHash(db *gorm.DB, u *models.User, plaintext string) error {
	if u == nil {
		return nil
	}
	if strings.HasPrefix(u.PasswordHash, "$argon2id$") {
		return nil
	}
	newHash, err := auth.HashPassword(plaintext)
	if err != nil {
		return err
	}
	return db.Model(&models.User{}).Where("id = ?", u.ID).Update("password_hash", newHash).Error
}
