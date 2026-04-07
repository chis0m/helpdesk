package mail

import "helpdesk/backend/internal/config"

// NewNotifiers returns staff-invite and password-reset notifiers based on MAIL_MAILER (log vs smtp).
func NewNotifiers(cfg config.Config) (StaffInviteNotifier, PasswordResetNotifier) {
	if cfg.UseSMTPMail() {
		m := NewSMTPMailer(cfg)
		return m, m
	}
	return NewLogStaffInviteNotifier(), NewLogPasswordResetNotifier()
}
