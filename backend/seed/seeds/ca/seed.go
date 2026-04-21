package ca

import (
	"fmt"

	"gorm.io/gorm"
)

func SeedAll(db *gorm.DB) error {
	staffSam, staffCassey, err := EnsureStaffUsers(db)
	if err != nil {
		return fmt.Errorf("ca seed: staff users: %w", err)
	}
	uMust, uMark, uJane, uAlex, err := EnsureCustomerUsers(db)
	if err != nil {
		return fmt.Errorf("ca seed: customer users: %w", err)
	}
	_ = uMust
	if err := EnsureTickets(db, uAlex, uMark, uJane, staffSam, staffCassey); err != nil {
		return fmt.Errorf("ca seed: tickets: %w", err)
	}
	if err := EnsureTicketComments(db, uMark, uJane, staffSam, staffCassey); err != nil {
		return fmt.Errorf("ca seed: ticket comments: %w", err)
	}
	return nil
}
