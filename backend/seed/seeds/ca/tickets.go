package ca

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"helpdesk/backend/internal/models"
)

// CA seed ticket titles (idempotent firstOrCreate by title). Must Change user has no tickets.
const (
	// John Doe (settled password) — three tickets
	TicketJohnOpenTitle       = "Sidebar keeps collapsing when I switch between projects"
	TicketJohnInProgressTitle = "CSV export times out after about 30 seconds"
	TicketJohnResolvedTitle   = "Cannot log in after password reset — now fixed on my side"
	// Jane Doe — two tickets
	TicketJaneOpenTitle  = "Question: when are monthly billing reminder emails sent?"
	TicketJaneClosedTitle = "Invoice PDF shows wrong VAT rate for last month's bill"
)

// EnsureTickets creates CA demo tickets: 3 for John, 2 for Jane. Assignees balanced between Sam (3) and Cassey (2).
func EnsureTickets(
	db *gorm.DB,
	uMust *models.User,
	uJohn *models.User,
	uJane *models.User,
	staffSam *models.User,
	staffCassey *models.User,
) error {
	_ = uMust // no tickets for first-login / must-change user

	if uJohn == nil || uJane == nil || staffSam == nil || staffCassey == nil {
		return fmt.Errorf("EnsureTickets: missing user(s)")
	}

	samID := staffSam.ID
	casseyID := staffCassey.ID

	// John — open (Sam)
	if _, err := firstOrCreateTicketWithStatus(
		db, uJohn.ID, TicketJohnOpenTitle,
		"Using the web app on Safari 17. The left sidebar collapses every time I click a different project in the dropdown. Expected it to stay open until I toggle it.",
		"general", models.TicketStatusOpen, &samID,
	); err != nil {
		return err
	}
	// John — in progress (Cassey)
	if _, err := firstOrCreateTicketWithStatus(
		db, uJohn.ID, TicketJohnInProgressTitle,
		"Exporting a medium-sized report (~5k rows) to CSV. The request runs then fails after roughly half a minute. Browser: Chrome on macOS.",
		"general", models.TicketStatusInProgress, &casseyID,
	); err != nil {
		return err
	}
	// John — resolved (Sam)
	if _, err := firstOrCreateTicketWithStatus(
		db, uJohn.ID, TicketJohnResolvedTitle,
		"After using the password reset link I was stuck in a loop on the login page. This started after the weekend maintenance window.",
		"general", models.TicketStatusResolved, &samID,
	); err != nil {
		return err
	}

	// Jane — open (Cassey)
	if _, err := firstOrCreateTicketWithStatus(
		db, uJane.ID, TicketJaneOpenTitle,
		"We need to coordinate with finance. Are billing reminders sent on the 1st, or on the invoice date? Thanks.",
		"billing", models.TicketStatusOpen, &casseyID,
	); err != nil {
		return err
	}
	// Jane — closed (Sam)
	if _, err := firstOrCreateTicketWithStatus(
		db, uJane.ID, TicketJaneClosedTitle,
		"The PDF attached to last month's billing email shows the wrong VAT percentage for our region. Need a corrected document or confirmation on the rate.",
		"billing", models.TicketStatusClosed, &samID,
	); err != nil {
		return err
	}

	return nil
}

func firstOrCreateTicketWithStatus(
	db *gorm.DB,
	reporterID uint64,
	title, description string,
	category string,
	status models.TicketStatus,
	assigneeID *uint64,
) (*models.Ticket, error) {
	normalizedTitle := strings.TrimSpace(title)
	var found []models.Ticket
	if err := db.Where("title = ?", normalizedTitle).Limit(1).Find(&found).Error; err != nil {
		return nil, err
	}
	if len(found) > 0 {
		return &found[0], nil
	}

	now := time.Now().UTC()
	t := models.Ticket{
		Title:          normalizedTitle,
		Description:    strings.TrimSpace(description),
		Status:         status,
		Category:       category,
		ReporterUserID: reporterID,
		AssignedUserID: assigneeID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if err := db.Create(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}
