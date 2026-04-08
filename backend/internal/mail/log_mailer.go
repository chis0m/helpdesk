package mail

import (
	"helpdesk/backend/internal/logger"
)

// Used when MAIL_MAILER=log (default): URLs are written to the app log (zerolog).
// Set MAIL_MAILER=smtp and MAIL_* variables to send via SMTP (e.g. Mailtrap).

type StaffInviteNotifier interface {
	SendStaffInvite(toEmail, inviteURL string) error
}

type LogStaffInviteNotifier struct{}

func NewLogStaffInviteNotifier() *LogStaffInviteNotifier {
	return &LogStaffInviteNotifier{}
}

func (n *LogStaffInviteNotifier) SendStaffInvite(toEmail, inviteURL string) error {
	logger.L().Info().
		Str("mail_driver", "log").
		Str("mail_outcome", "success").
		Str("kind", "staff_invite").
		Str("to", toEmail).
		Str("invite_url", inviteURL).
		Msg("mail: staff invite written to logs successfully (not sent via SMTP)")
	return nil
}

// Password reset: same as staff invite — log when MAIL_MAILER=log.

type PasswordResetNotifier interface {
	SendPasswordReset(toEmail, resetURL string) error
}

type LogPasswordResetNotifier struct{}

func NewLogPasswordResetNotifier() *LogPasswordResetNotifier {
	return &LogPasswordResetNotifier{}
}

func (n *LogPasswordResetNotifier) SendPasswordReset(toEmail, resetURL string) error {
	logger.L().Info().
		Str("mail_driver", "log").
		Str("mail_outcome", "success").
		Str("kind", "password_reset").
		Str("to", toEmail).
		Str("reset_url", resetURL).
		Msg("mail: password reset link written to logs successfully (not sent via SMTP)")
	return nil
}
