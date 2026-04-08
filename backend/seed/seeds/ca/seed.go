package ca

import (
	"fmt"

	"gorm.io/gorm"
)

// SeedAll applies the CA assessment seed: staff + three customers + five tickets (Mark 3, Jane 2; Must Change has none).
func SeedAll(db *gorm.DB) error {
	staffSam, staffCassey, err := EnsureStaffUsers(db)
	if err != nil {
		return fmt.Errorf("ca seed: staff users: %w", err)
	}
	uMust, uMark, uJane, err := EnsureCustomerUsers(db)
	if err != nil {
		return fmt.Errorf("ca seed: customer users: %w", err)
	}
	if err := EnsureTickets(db, uMust, uMark, uJane, staffSam, staffCassey); err != nil {
		return fmt.Errorf("ca seed: tickets: %w", err)
	}
	if err := EnsureTicketComments(db, uMark, uJane, staffSam, staffCassey); err != nil {
		return fmt.Errorf("ca seed: ticket comments: %w", err)
	}
	return nil
}
