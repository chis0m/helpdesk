package ca

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"helpdesk/backend/internal/auth"
	"helpdesk/backend/internal/models"
)

// Exported emails for TASK.md / boot docs (CA seed customer users).
const (
	EmailMustChange = "must.change@example.com"
	EmailJohnDoe    = "john.doe@sample.com"
	EmailJaneDoe    = "jane.doe@sample.com"
)

// EnsureCustomerUsers creates three CA demo customers: Must Change (first login), John Doe, Jane Doe (both passwords settled).
// Password hashes use Argon2id (same as auth login). Legacy bcrypt rows are upgraded on re-seed.
func EnsureCustomerUsers(db *gorm.DB) (*models.User, *models.User, *models.User, error) {
	hash, err := auth.HashPassword(caTestPassword)
	if err != nil {
		return nil, nil, nil, err
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
		return nil, nil, nil, err
	}
	if err := ensureArgon2idPasswordHash(db, uMust, caTestPassword); err != nil {
		return nil, nil, nil, err
	}

	uJohn, _, err := firstOrCreateUser(
		db,
		EmailJohnDoe,
		"John",
		"Doe",
		models.RoleUser,
		false,
		&settled,
		hash,
	)
	if err != nil {
		return nil, nil, nil, err
	}
	if err := ensureArgon2idPasswordHash(db, uJohn, caTestPassword); err != nil {
		return nil, nil, nil, err
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
		return nil, nil, nil, err
	}
	if err := ensureArgon2idPasswordHash(db, uJane, caTestPassword); err != nil {
		return nil, nil, nil, err
	}

	if uMust == nil || uJohn == nil || uJane == nil {
		return nil, nil, nil, fmt.Errorf("EnsureCustomerUsers: missing user(s)")
	}
	return uMust, uJohn, uJane, nil
}

// ensureArgon2idPasswordHash re-hashes with Argon2id if the stored hash is legacy (e.g. bcrypt from an older seed).
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
