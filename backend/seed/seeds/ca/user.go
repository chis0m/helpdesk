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
	customerDomainRileyOrg = "acmelogistics.ie"
	customerDomainJordanOrg = "northwind.ie"
)

// EmailCustomerMustChange is the CA user who must change password on first login (known password).
var EmailCustomerMustChange = emailAt("Riley", "MustChange", customerDomainRileyOrg)

// EmailCustomerOK is the CA user with password already “settled” (random password logged on first create).
var EmailCustomerOK = emailAt("Jordan", "Lee", customerDomainJordanOrg)

// hardcodedMustChangeUserPassword is the only CA credential stored in code — for testing forced password change.
// It is logged once on first seed create; in production or shared environments this value MUST be changed / not reused.
const hardcodedMustChangeUserPassword = "CaMustChange1!"

// EnsureCustomerUsers seeds two portal users. Riley uses hardcodedMustChangeUserPassword; Jordan gets a random password (logged on first creation only).
func EnsureCustomerUsers(db *gorm.DB) (mustChangeUser *models.User, okUser *models.User, err error) {
	log := logger.L()
	now := time.Now().UTC()
	ptrNow := &now

	hashMust, err := auth.HashPassword(hardcodedMustChangeUserPassword)
	if err != nil {
		return nil, nil, err
	}
	uMust, createdMust, err := firstOrCreateUser(
		db,
		EmailCustomerMustChange,
		"Riley",
		"MustChange",
		models.RoleUser,
		true,
		nil,
		hashMust,
	)
	if err != nil {
		return nil, nil, err
	}
	if createdMust {
		// Log plaintext so assessors can sign in once. This row has must_change_password=true — the real password MUST be changed on first login (do not treat as a long-term secret).
		log.Warn().
			Str("email", uMust.Email).
			Str("plaintext_password", hardcodedMustChangeUserPassword).
			Msg("CA seed: initial test password logged on purpose — this account MUST change password on first login (intentional for CA; not shown again if user already exists)")
	}

	plainRand, err := randomPassword()
	if err != nil {
		return nil, nil, err
	}
	hashOK, err := auth.HashPassword(plainRand)
	if err != nil {
		return nil, nil, err
	}
	uOK, createdOK, err := firstOrCreateUser(
		db,
		EmailCustomerOK,
		"Jordan",
		"Lee",
		models.RoleUser,
		false,
		ptrNow,
		hashOK,
	)
	if err != nil {
		return nil, nil, err
	}
	if createdOK {
		log.Warn().
			Str("email", uOK.Email).
			Str("plaintext_password", plainRand).
			Msg("CA seed: random password generated on purpose for test data — save this to sign in (intentional; not shown again if user already exists)")
	}

	return uMust, uOK, nil
}
