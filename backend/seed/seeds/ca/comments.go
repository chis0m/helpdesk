package ca

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"helpdesk/backend/internal/models"
)

// EnsureTicketComments adds realistic thread comments: John's in_progress + resolved tickets; Jane's closed ticket (5).
// Staff authors are Sam Support and Cassey Support. Idempotent by (ticket_id, body).
func EnsureTicketComments(
	db *gorm.DB,
	uJohn *models.User,
	uJane *models.User,
	staffSam *models.User,
	staffCassey *models.User,
) error {
	if uJohn == nil || uJane == nil || staffSam == nil || staffCassey == nil {
		return fmt.Errorf("EnsureTicketComments: missing user(s)")
	}

	tInProg, err := ticketByTitle(db, TicketJohnInProgressTitle)
	if err != nil {
		return err
	}
	tResolved, err := ticketByTitle(db, TicketJohnResolvedTitle)
	if err != nil {
		return err
	}
	tJaneClosed, err := ticketByTitle(db, TicketJaneClosedTitle)
	if err != nil {
		return err
	}

	samID := staffSam.ID
	casseyID := staffCassey.ID
	johnID := uJohn.ID
	janeID := uJane.ID

	// John — in progress (Cassey's ticket): dialogue with Cassey
	base := time.Now().UTC().Add(-72 * time.Hour)
	if err := seedCommentChain(db, tInProg.ID, []commentSeed{
		{johnID, "Hi — I'm still seeing the timeout after about 30s on CSV export. Happy to send logs if useful.", base.Add(0 * time.Minute)},
		{casseyID, "Thanks, John. Could you grab a HAR from the browser devtools for the failing export request and confirm your browser version?", base.Add(18 * time.Minute)},
		{johnID, "Attached — failing request is POST /api/export. Chrome 131 on macOS.", base.Add(42 * time.Minute)},
		{casseyID, "Received — we can reproduce internally. Engineering is looking at it; I'll post an update when we have a fix.", base.Add(3 * time.Hour)},
	}); err != nil {
		return err
	}

	// John — resolved (Sam's ticket): John, Sam, Cassey
	base2 := time.Now().UTC().Add(-120 * time.Hour)
	if err := seedCommentChain(db, tResolved.ID, []commentSeed{
		{johnID, "Quick update: I can log in normally again after the patch you deployed on Tuesday.", base2.Add(0 * time.Minute)},
		{samID, "Thanks for confirming — I'll mark this resolved from support.", base2.Add(25 * time.Minute)},
		{casseyID, "Noted for our weekly review — appreciate the fast turnaround, Sam and John.", base2.Add(50 * time.Minute)},
	}); err != nil {
		return err
	}

	// Jane — closed: five-comment history (Jane, Sam, Jane, Cassey, Jane)
	base3 := time.Now().UTC().Add(-200 * time.Hour)
	if err := seedCommentChain(db, tJaneClosed.ID, []commentSeed{
		{janeID, "The invoice PDF in last month's billing email still shows the wrong VAT rate for our region.", base3.Add(0 * time.Minute)},
		{samID, "Thanks — please share the invoice number and roughly when you received the email.", base3.Add(2 * time.Hour)},
		{janeID, "Invoice INV-2024-0099 — received the evening of 18 March.", base3.Add(5 * time.Hour)},
		{casseyID, "Found it on our side — we've corrected the VAT line. Please re-download the PDF from the billing portal (same link).", base3.Add(24 * time.Hour)},
		{janeID, "Re-downloaded just now — the PDF matches our records. Thanks, we can consider this closed.", base3.Add(26 * time.Hour)},
	}); err != nil {
		return err
	}

	return nil
}

type commentSeed struct {
	authorID  uint64
	body      string
	createdAt time.Time
}

func seedCommentChain(db *gorm.DB, ticketID uint64, chain []commentSeed) error {
	for _, c := range chain {
		if err := firstOrCreateComment(db, ticketID, c.authorID, c.body, c.createdAt); err != nil {
			return err
		}
	}
	return nil
}

func ticketByTitle(db *gorm.DB, title string) (*models.Ticket, error) {
	var rows []models.Ticket
	t := strings.TrimSpace(title)
	if err := db.Where("title = ?", t).Limit(1).Find(&rows).Error; err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("ticket not found for title %q (run ticket seed first)", t)
	}
	return &rows[0], nil
}

func firstOrCreateComment(db *gorm.DB, ticketID, authorUserID uint64, body string, createdAt time.Time) error {
	b := strings.TrimSpace(body)
	var found []models.TicketComment
	if err := db.Where("ticket_id = ? AND body = ?", ticketID, b).Limit(1).Find(&found).Error; err != nil {
		return err
	}
	if len(found) > 0 {
		return nil
	}

	c := models.TicketComment{
		TicketID:     ticketID,
		AuthorUserID: authorUserID,
		Body:         b,
		CreatedAt:    createdAt,
		UpdatedAt:    createdAt,
	}
	return db.Create(&c).Error
}
