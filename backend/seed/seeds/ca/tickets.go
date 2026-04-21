package ca

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"helpdesk/backend/internal/models"
)

const (
	TicketMarkOpenTitle        = "Sidebar keeps collapsing when I switch between projects"
	TicketMarkInProgressTitle  = "CSV export times out after about 30 seconds"
	TicketMarkResolvedTitle    = "Cannot log in after password reset — now fixed on my side"
	TicketJaneOpenTitle        = "Question: how do charges for SecWeb technical services appear on our monthly bill?"
	TicketJaneOpen2Title       = "Dashboard summary numbers look out of date until I force a full page reload"
	TicketJaneVPNTitle         = "Can't keep the VPN connected today. I'm out for meetings all day and need remote access to work."
	TicketTechSam1Title        = "M365 shared mailbox: calendar free/busy not updating after SecWeb migration cutover"
	TicketTechSam2Title        = "DNS: delegated subdomain for customer portal still resolving to old IP after change request"
	TicketTechSam3Title        = "Endpoint: BitLocker recovery key prompt loop on fleet laptops after SecWeb policy bundle"
	TicketTechCassey1Title     = "Firewall: need rule review for new outbound SFTP to vendor maintenance subnet"
	TicketTechCassey2Title     = "Site Wi-Fi: conference wing drops every few minutes — request heatmap and AP tuning"
	TicketTechCassey3Title     = "SSO/SAML: assertion failures from IdP after certificate rotation — SecWeb SP metadata check"
	TicketTechUnassigned1Title = "Server monitoring: Hyper-V host CPU pegged post-patch — need SecWeb infra triage window"
	TicketTechUnassigned2Title = "Managed print queue on print01 not draining jobs — affects finance batch printing"
	TicketTechUnassigned3Title = "Backup verification job failing on SQL maintenance plan — SecWeb BCDR checklist item"
)

func EnsureTickets(
	db *gorm.DB,
	uAlex *models.User,
	uMark *models.User,
	uJane *models.User,
	staffSam *models.User,
	staffCassey *models.User,
) error {
	if uMark == nil || uJane == nil || uAlex == nil || staffSam == nil || staffCassey == nil {
		return fmt.Errorf("EnsureTickets: missing user(s)")
	}

	samID := staffSam.ID
	casseyID := staffCassey.ID

	if _, err := firstOrCreateTicketWithStatus(
		db, uMark.ID, TicketMarkOpenTitle,
		"Using the web app on Safari 17. The left sidebar collapses every time I click a different project in the dropdown. Expected it to stay open until I toggle it.",
		"general", models.TicketStatusOpen, &samID,
	); err != nil {
		return err
	}
	if _, err := firstOrCreateTicketWithStatus(
		db, uMark.ID, TicketMarkInProgressTitle,
		"Exporting a medium-sized report (~5k rows) to CSV. The request runs then fails after roughly half a minute. Browser: Chrome on macOS.",
		"general", models.TicketStatusInProgress, &casseyID,
	); err != nil {
		return err
	}
	if _, err := firstOrCreateTicketWithStatus(
		db, uMark.ID, TicketMarkResolvedTitle,
		"After using the password reset link I was stuck in a loop on the login page. This started after the weekend maintenance window.",
		"general", models.TicketStatusResolved, &samID,
	); err != nil {
		return err
	}

	if _, err := firstOrCreateTicketWithStatus(
		db, uJane.ID, TicketJaneOpenTitle,
		"Our finance team is reviewing last quarter. I need to understand how fees for technical services SecWeb provides (extra support hours, onboarding assistance, scoped projects) appear on the invoice alongside the standard monthly charge. Details in the comments.",
		"billing", models.TicketStatusOpen, &casseyID,
	); err != nil {
		return err
	}
	if _, err := firstOrCreateTicketWithStatus(
		db, uJane.ID, TicketJaneOpen2Title,
		"The KPI tiles on my SecWeb dashboard often show yesterday’s figures until I do a hard refresh (Cmd+Shift+R). Colleagues see the same. Browser: Edge on Windows 11.",
		"general", models.TicketStatusOpen, &samID,
	); err != nil {
		return err
	}
	if _, err := firstOrCreateTicketWithStatus(
		db, uJane.ID, TicketJaneVPNTitle,
		"I'm unable to stay connected to the company VPN while I'm out for meetings today. I need remote access for work — more detail in the comments below.",
		"general", models.TicketStatusInProgress, &samID,
	); err != nil {
		return err
	}
	if err := db.Model(&models.Ticket{}).
		Where("title = ?", strings.TrimSpace(TicketJaneVPNTitle)).
		Update("status", models.TicketStatusInProgress).Error; err != nil {
		return err
	}

	if _, err := firstOrCreateTicketWithStatus(
		db, uMark.ID, TicketTechSam1Title,
		"Finance shared mailbox migrated Friday. Internal users see stale free/busy in Outlook. SecWeb handles our tenant — need engineer to verify connector and replication lag.",
		"technical", models.TicketStatusOpen, &samID,
	); err != nil {
		return err
	}
	if _, err := firstOrCreateTicketWithStatus(
		db, uJane.ID, TicketTechSam2Title,
		"Submitted DNS change ticket #4412 through the portal. TTL lowered to 300. NS at registrar points to SecWeb DNS but A record still shows previous hosting IP after 24h.",
		"technical", models.TicketStatusInProgress, &samID,
	); err != nil {
		return err
	}
	if _, err := firstOrCreateTicketWithStatus(
		db, uAlex.ID, TicketTechSam3Title,
		"Several laptops from the rollout group prompt for BitLocker recovery on every cold boot since yesterday’s policy sync. SecWeb MDM shows compliant — need root cause.",
		"technical", models.TicketStatusOpen, &samID,
	); err != nil {
		return err
	}
	if _, err := firstOrCreateTicketWithStatus(
		db, uJane.ID, TicketTechCassey1Title,
		"Vendor requires outbound SFTP to /24 range for nightly extracts. Security wants SecWeb to validate rule scope and logging before go-live.",
		"technical", models.TicketStatusOpen, &casseyID,
	); err != nil {
		return err
	}
	if _, err := firstOrCreateTicketWithStatus(
		db, uMark.ID, TicketTechCassey2Title,
		"Users in rooms 201–205 lose Wi-Fi under load. SecWeb manages our Aruba stack — need onsite or remote survey and channel/power recommendations.",
		"technical", models.TicketStatusInProgress, &casseyID,
	); err != nil {
		return err
	}
	if _, err := firstOrCreateTicketWithStatus(
		db, uAlex.ID, TicketTechCassey3Title,
		"IdP rotated signing cert Sunday. SAML logins to SecWeb-managed apps return invalid signature. Metadata upload may be out of date.",
		"technical", models.TicketStatusOpen, &casseyID,
	); err != nil {
		return err
	}
	if _, err := firstOrCreateTicketWithStatus(
		db, uMark.ID, TicketTechUnassigned1Title,
		"HV-CL02 running customer workloads. CPU sustained >90% since Feb cumulative updates. Need SecWeb infra team to review guest placement and host patches.",
		"technical", models.TicketStatusOpen, nil,
	); err != nil {
		return err
	}
	if _, err := firstOrCreateTicketWithStatus(
		db, uJane.ID, TicketTechUnassigned2Title,
		"Print server print01 under SecWeb managed services. Jobs stuck in queue; driver refresh attempted locally without success.",
		"technical", models.TicketStatusOpen, nil,
	); err != nil {
		return err
	}
	if _, err := firstOrCreateTicketWithStatus(
		db, uAlex.ID, TicketTechUnassigned3Title,
		"Weekly verify-only backup job errors on instance SQL-REP-03. Impacts our BCDR evidence pack for auditors next month.",
		"technical", models.TicketStatusInProgress, nil,
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
