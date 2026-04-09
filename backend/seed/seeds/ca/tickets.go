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
	// Mark Anthony (settled password) — three tickets
	TicketMarkOpenTitle       = "Sidebar keeps collapsing when I switch between projects"
	TicketMarkInProgressTitle = "CSV export times out after about 30 seconds"
	TicketMarkResolvedTitle   = "Cannot log in after password reset — now fixed on my side"
	// Jane Doe — two tickets
	TicketJaneOpenTitle = "Question: how do charges for SecWeb technical services appear on our monthly bill?"
	TicketJaneClosedTitle = "Can't keep the VPN connected today. I'm out for meetings all day and need remote access to work."
)

// EnsureTickets creates CA demo tickets: 3 for Mark, 2 for Jane. Assignees balanced between Sam (3) and Cassey (2).
func EnsureTickets(
	db *gorm.DB,
	uMust *models.User,
	uMark *models.User,
	uJane *models.User,
	staffSam *models.User,
	staffCassey *models.User,
) error {
	_ = uMust // no tickets for first-login / must-change user

	if uMark == nil || uJane == nil || staffSam == nil || staffCassey == nil {
		return fmt.Errorf("EnsureTickets: missing user(s)")
	}

	samID := staffSam.ID
	casseyID := staffCassey.ID

	// Mark — open (Sam)
	if _, err := firstOrCreateTicketWithStatus(
		db, uMark.ID, TicketMarkOpenTitle,
		"Using the web app on Safari 17. The left sidebar collapses every time I click a different project in the dropdown. Expected it to stay open until I toggle it.",
		"general", models.TicketStatusOpen, &samID,
	); err != nil {
		return err
	}
	// Mark — in progress (Cassey)
	if _, err := firstOrCreateTicketWithStatus(
		db, uMark.ID, TicketMarkInProgressTitle,
		"Exporting a medium-sized report (~5k rows) to CSV. The request runs then fails after roughly half a minute. Browser: Chrome on macOS.",
		"general", models.TicketStatusInProgress, &casseyID,
	); err != nil {
		return err
	}
	// Mark — resolved (Sam)
	if _, err := firstOrCreateTicketWithStatus(
		db, uMark.ID, TicketMarkResolvedTitle,
		"After using the password reset link I was stuck in a loop on the login page. This started after the weekend maintenance window.",
		"general", models.TicketStatusResolved, &samID,
	); err != nil {
		return err
	}

	// Jane — open (Cassey): billing for technical / professional services vs. base subscription.
	if _, err := firstOrCreateTicketWithStatus(
		db, uJane.ID, TicketJaneOpenTitle,
		"Our finance team is reviewing last quarter. I need to understand how fees for technical services SecWeb provides (extra support hours, onboarding assistance, scoped projects) appear on the invoice alongside the standard monthly charge. Details in the comments.",
		"billing", models.TicketStatusOpen, &casseyID,
	); err != nil {
		return err
	}
	// Jane — closed (Sam): VPN thread; see ca/comments.go for credential-in-ticket demo (IDOR impact).
	if _, err := firstOrCreateTicketWithStatus(
		db, uJane.ID, TicketJaneClosedTitle,
		"I'm unable to stay connected to the company VPN while I'm out for meetings today. I need remote access for work — more detail in the comments below.",
		"general", models.TicketStatusClosed, &samID,
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
