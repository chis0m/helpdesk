package ca

import (
	"gorm.io/gorm"

	"helpdesk/backend/internal/logger"
	"helpdesk/backend/internal/models"
)

// Ticket titles (idempotent).
const (
	TicketNoCommentsTitle    = "[CA Seed] Unassigned — no comments"
	TicketThreeCommentsTitle = "[CA Seed] Assigned — 3-comment thread"
	TicketUserOnlyTitle      = "[CA Seed] Unassigned — user comment only"
)

// EnsureTickets seeds three CA tickets when missing (by title).
func EnsureTickets(db *gorm.DB, uMust, uOK, staffSupport *models.User) error {
	if err := firstOrCreateTicketNoComments(db, uMust.ID, TicketNoCommentsTitle); err != nil {
		return err
	}
	if err := firstOrCreateTicketThreeComments(db, uOK.ID, staffSupport.ID, TicketThreeCommentsTitle, uOK.ID, staffSupport.ID); err != nil {
		return err
	}
	if err := firstOrCreateTicketUserCommentOnly(db, uMust.ID, TicketUserOnlyTitle, uMust.ID); err != nil {
		return err
	}
	return nil
}

func firstOrCreateTicketNoComments(db *gorm.DB, reporterID uint64, title string) error {
	var existing []models.Ticket
	if err := db.Where("title = ?", title).Limit(1).Find(&existing).Error; err != nil {
		return err
	}
	if len(existing) > 0 {
		return nil
	}

	t := models.Ticket{
		ReporterUserID: reporterID,
		AssignedUserID: nil,
		Title:          title,
		Description:    "CA seed ticket: no comments, unassigned. Reporter is the user who must change password on first login.",
		Category:       "general",
		Status:         models.TicketStatusOpen,
	}
	if err := db.Create(&t).Error; err != nil {
		return err
	}
	logger.L().Info().Str("title", title).Msg("CA fixture ticket created")
	return nil
}

func firstOrCreateTicketThreeComments(db *gorm.DB, reporterID, assigneeID uint64, title string, authorUserID, authorStaffID uint64) error {
	var existing []models.Ticket
	if err := db.Where("title = ?", title).Limit(1).Find(&existing).Error; err != nil {
		return err
	}
	if len(existing) > 0 {
		return nil
	}

	assignee := assigneeID
	t := models.Ticket{
		ReporterUserID: reporterID,
		AssignedUserID: &assignee,
		Title:          title,
		Description:    "CA seed ticket: three comments (reporter + staff + reporter), assigned to support staff.",
		Category:       "technical",
		Status:         models.TicketStatusInProgress,
	}
	if err := db.Create(&t).Error; err != nil {
		return err
	}

	comments := []struct {
		authorID uint64
		body     string
	}{
		{authorUserID, "Initial report: intermittent sync failure after the last update."},
		{authorStaffID, "Thanks — I can reproduce. Assigned to me; investigating logs."},
		{authorUserID, "Follow-up: I attached export-logs.zip to the case."},
	}
	for _, c := range comments {
		co := models.TicketComment{
			TicketID:     t.ID,
			AuthorUserID: c.authorID,
			Body:         c.body,
		}
		if err := db.Create(&co).Error; err != nil {
			return err
		}
	}

	logger.L().Info().Str("title", title).Msg("CA fixture ticket + comments created")
	return nil
}

func firstOrCreateTicketUserCommentOnly(db *gorm.DB, reporterID uint64, title string, commentAuthorID uint64) error {
	var existing []models.Ticket
	if err := db.Where("title = ?", title).Limit(1).Find(&existing).Error; err != nil {
		return err
	}
	if len(existing) > 0 {
		return nil
	}

	t := models.Ticket{
		ReporterUserID: reporterID,
		AssignedUserID: nil,
		Title:          title,
		Description:    "CA seed ticket: single comment from the reporter only, unassigned.",
		Category:       "billing",
		Status:         models.TicketStatusOpen,
	}
	if err := db.Create(&t).Error; err != nil {
		return err
	}

	co := models.TicketComment{
		TicketID:     t.ID,
		AuthorUserID: commentAuthorID,
		Body:         "Question on invoice line items — can someone confirm tax is correct?",
	}
	if err := db.Create(&co).Error; err != nil {
		return err
	}

	logger.L().Info().Str("title", title).Msg("CA fixture ticket + user comment created")
	return nil
}
