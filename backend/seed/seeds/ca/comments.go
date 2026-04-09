package ca

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"helpdesk/backend/internal/models"
)

// EnsureTicketComments adds realistic thread comments: Mark's in_progress + resolved; Jane's open (3); Jane's closed (4).
// Staff authors are Sam Support and Cassey Support. Idempotent by (ticket_id, body).
func EnsureTicketComments(
	db *gorm.DB,
	uMark *models.User,
	uJane *models.User,
	staffSam *models.User,
	staffCassey *models.User,
) error {
	if uMark == nil || uJane == nil || staffSam == nil || staffCassey == nil {
		return fmt.Errorf("EnsureTicketComments: missing user(s)")
	}

	tInProg, err := ticketByTitle(db, TicketMarkInProgressTitle)
	if err != nil {
		return err
	}
	tResolved, err := ticketByTitle(db, TicketMarkResolvedTitle)
	if err != nil {
		return err
	}
	tJaneOpen, err := ticketByTitle(db, TicketJaneOpenTitle)
	if err != nil {
		return err
	}
	tJaneClosed, err := ticketByTitle(db, TicketJaneClosedTitle)
	if err != nil {
		return err
	}

	samID := staffSam.ID
	casseyID := staffCassey.ID
	markID := uMark.ID
	janeID := uJane.ID

	// Mark — in progress (Cassey's ticket): dialogue with Cassey
	base := time.Now().UTC().Add(-72 * time.Hour)
	if err := seedCommentChain(db, tInProg.ID, []commentSeed{
		{markID, "Hi — I'm still seeing the timeout after about 30s on CSV export. Happy to send logs if useful.", base.Add(0 * time.Minute)},
		{casseyID, "Thanks, Mark. Could you grab a HAR from the browser devtools for the failing export request and confirm your browser version?", base.Add(18 * time.Minute)},
		{markID, "Attached — failing request is POST /api/export. Chrome 131 on macOS.", base.Add(42 * time.Minute)},
		{casseyID, "Received — we can reproduce internally. Engineering is looking at it; I'll post an update when we have a fix.", base.Add(3 * time.Hour)},
	}); err != nil {
		return err
	}

	// Mark — resolved (Sam's ticket): Mark, Sam, Cassey
	base2 := time.Now().UTC().Add(-120 * time.Hour)
	if err := seedCommentChain(db, tResolved.ID, []commentSeed{
		{markID, "Quick update: I can log in normally again after the patch you deployed on Tuesday.", base2.Add(0 * time.Minute)},
		{samID, "Thanks for confirming — I'll mark this resolved from support.", base2.Add(25 * time.Minute)},
		{casseyID, "Noted for our weekly review — appreciate the fast turnaround, Sam and Mark.", base2.Add(50 * time.Minute)},
	}); err != nil {
		return err
	}

	// Jane — open: technical / professional services billing vs. base subscription (Cassey assigned).
	baseJaneOpen := time.Now().UTC().Add(-150 * time.Hour)
	if err := seedCommentChain(db, tJaneOpen.ID, []commentSeed{
		{janeID, "We're consolidating costs with finance. Can you confirm which invoice lines cover technical services SecWeb has delivered for us — extra support hours, the onboarding work, anything beyond the standard monthly fee? I need to avoid double-counting with our internal numbers.", baseJaneOpen},
		{casseyID, "Hi Jane — the base subscription usually shows as one line item. Technical or professional services (additional hours, onboarding or one-off projects, anything scoped outside the plan) are listed separately with their own descriptions and rates. I can ask billing to email you a sample invoice layout from your account if that helps.", baseJaneOpen.Add(50 * time.Minute)},
		{janeID, "That's what we needed — yes, please have billing send a sample. I'll align with finance once we have it.", baseJaneOpen.Add(3 * time.Hour)},
	}); err != nil {
		return err
	}

	// Jane — closed: VPN + pasted credentials; support asks removal; Jane can't get online to edit; Cassey gentle reminder 4h later (IDOR demo).
	base3 := time.Now().UTC().Add(-200 * time.Hour)
	janeFollowUp := base3.Add(3 * time.Hour)
	if err := seedCommentChain(db, tJaneClosed.ID, []commentSeed{
		{janeID, "Can't keep the VPN connected today — I'm in meetings all day and need this working. I've tried logging in several times and it keeps dropping me out. I'm listed as an admin on our company account if that helps.\n\nI'm pasting what I use in case it's on our side — please don't share this:\nUsername: jane.doe.acme\nPassword: FakeVPN-Seed-2024-NotReal\n\nWon't be at a real keyboard until late — please fix if you can.", base3},
		{samID, "Sorry you're stuck on this. Please delete your previous message — don't put passwords in the ticket. Use a password manager or secure channel next time. We're checking VPN on our end.", base3.Add(1 * time.Hour)},
		{janeID, "Got it. Hotel Wi-Fi is awful — I can't get the portal to load reliably to edit that message right now. I'll remove it when I'm on something stable.", janeFollowUp},
		{casseyID, "Hi Jane — gentle reminder when you have a moment: could you remove the username and password from your first message so they aren't kept in the ticket? Thanks.", janeFollowUp.Add(4 * time.Hour)},
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
