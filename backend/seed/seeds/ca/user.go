package ca

import (
	"time"

	"gorm.io/gorm"

	"helpdesk/backend/internal/auth"
	"helpdesk/backend/internal/logger"
	"helpdesk/backend/internal/models"
)

// Customer org domains (firstname.lastname@company).
const (
	customerDomainRileyOrg  = "acmelogistics.ie"
	customerDomainJordanOrg = "northwind.ie"
)

// EmailCustomerMustChange is the CA user who must change password on first login.
var EmailCustomerMustChange = emailAt("Riley", "MustChange", customerDomainRileyOrg)

// EmailCustomerOK is the CA user with password already “settled” (same test password as others; PasswordChangedAt set).
var EmailCustomerOK = emailAt("Jordan", "Lee", customerDomainJordanOrg)

// EnsureCustomerUsers seeds two portal users. All use caTestPassword; Riley has MustChangePassword true, Jordan false.
func EnsureCustomerUsers(db *gorm.DB) (mustChangeUser *models.User, okUser *models.User, err error) {
	now := time.Now().UTC()
	ptrNow := &now

	hash, err := auth.HashPassword(caTestPassword)
	if err != nil {
		return nil, nil, err
	}

	log := logger.L()

	uMust, createdMust, err := firstOrCreateUser(
		db,
		EmailCustomerMustChange,
		"Riley",
		"MustChange",
		models.RoleUser,
		true,
		nil,
		hash,
	)
	if err != nil {
		return nil, nil, err
	}
	if createdMust {
		log.Info().
			Str("email", uMust.Email).
			Bool("must_change_password", uMust.MustChangePassword).
			Msg("CA seed: test user created")
	}

	uOK, createdOK, err := firstOrCreateUser(
		db,
		EmailCustomerOK,
		"Jordan",
		"Lee",
		models.RoleUser,
		false,
		ptrNow,
		hash,
	)
	if err != nil {
		return nil, nil, err
	}
	if createdOK {
		log.Info().
			Str("email", uOK.Email).
			Bool("must_change_password", uOK.MustChangePassword).
			Msg("CA seed: test user created")
	}

	return uMust, uOK, nil
}
