package mail

import (
	"helpdesk/backend/internal/logger"
)

// For the CA we only log the invite URL to the console so it can be copied from logs.
// In production this would be replaced with SMTP or a transactional email provider.

type StaffInviteNotifier interface {
	SendStaffInvite(toEmail, inviteURL string) error
}

type LogStaffInviteNotifier struct{}

func NewLogStaffInviteNotifier() *LogStaffInviteNotifier {
	return &LogStaffInviteNotifier{}
}

func (n *LogStaffInviteNotifier) SendStaffInvite(toEmail, inviteURL string) error {
	logger.L().Info().
		Str("to", toEmail).
		Str("invite_url", inviteURL).
		Msg("CA: staff invite link (email not sent; copy URL from logs)")
	return nil
}

// Password reset uses the same CA pattern as staff invites: log the URL, no SMTP.

type PasswordResetNotifier interface {
	SendPasswordReset(toEmail, resetURL string) error
}

type LogPasswordResetNotifier struct{}

func NewLogPasswordResetNotifier() *LogPasswordResetNotifier {
	return &LogPasswordResetNotifier{}
}

func (n *LogPasswordResetNotifier) SendPasswordReset(toEmail, resetURL string) error {
	logger.L().Info().
		Str("to", toEmail).
		Str("reset_url", resetURL).
		Msg("CA: password reset link (email not sent; copy URL from logs)")
	return nil
}
