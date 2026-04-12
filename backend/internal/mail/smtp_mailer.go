package mail

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"net"
	"net/smtp"
	"strconv"
	"strings"

	"helpdesk/backend/internal/config"
	"helpdesk/backend/internal/logger"
)

// SMTPMailer sends transactional mail via SMTP (e.g. Mailtrap). Implements StaffInviteNotifier and PasswordResetNotifier.
type SMTPMailer struct {
	cfg config.Config
}

// NewSMTPMailer returns an SMTP-backed mailer. Use only when config.UseSMTPMail() is true.
func NewSMTPMailer(cfg config.Config) *SMTPMailer {
	return &SMTPMailer{cfg: cfg}
}

func (s *SMTPMailer) SendStaffInvite(toEmail, inviteURL string) error {
	if err := s.validate(); err != nil {
		logger.L().Error().Err(err).
			Str("mail_driver", "smtp").
			Str("mail_outcome", "failed").
			Str("kind", "staff_invite").
			Str("to", toEmail).
			Msg("mail: smtp configuration validation failed")
		return err
	}
	app := strings.TrimSpace(s.cfg.AppName)
	if app == "" {
		app = "SecWeb HelpDesk"
	}
	subject := "Staff invitation"
	plain := fmt.Sprintf(
		"You’re invited to join the team.\n\nAccept and set your password:\n%s\n",
		inviteURL,
	)
	bodyHTML, err := renderEmailLayout(EmailLayout{
		AppName:     app,
		Badge:       "Staff invite",
		BodyIntro:   template.HTML(`<strong>You’re invited</strong> to join the team. Use the button below to accept this invitation and set your password.`),
		BodyMuted:   "This link expires in 15 minutes.",
		ButtonLabel: "Accept invitation",
		ButtonURL:   inviteURL,
		LinkHint:    "If the button doesn’t work, paste this URL into your browser:",
		LinkURL:     inviteURL,
		Footer:      "You received this because an administrator invited this email address. If this wasn’t you, you can ignore this message.",
	})
	if err != nil {
		logger.L().Error().Err(err).
			Str("mail_driver", "smtp").
			Str("mail_outcome", "failed").
			Str("kind", "staff_invite").
			Str("to", toEmail).
			Msg("mail: smtp email HTML render failed")
		return err
	}
	return s.send(toEmail, subject, plain, bodyHTML, "staff_invite")
}

func (s *SMTPMailer) SendPasswordReset(toEmail, resetURL string) error {
	if err := s.validate(); err != nil {
		logger.L().Error().Err(err).
			Str("mail_driver", "smtp").
			Str("mail_outcome", "failed").
			Str("kind", "password_reset").
			Str("to", toEmail).
			Msg("mail: smtp configuration validation failed")
		return err
	}
	app := strings.TrimSpace(s.cfg.AppName)
	if app == "" {
		app = "SecWeb HelpDesk"
	}
	subject := "Password reset"
	plain := fmt.Sprintf(
		"We received a request to reset your password.\n\nOpen this link:\n%s\n\nIf you did not request this, ignore this email.\n",
		resetURL,
	)
	bodyHTML, err := renderEmailLayout(EmailLayout{
		AppName:     app,
		Badge:       "Password reset",
		BodyIntro:   template.HTML(`We received a request to <strong>reset your password</strong>. Click the button to choose a new password.`),
		BodyMuted:   "If you didn’t ask for this, no action is needed — your password will stay the same.",
		ButtonLabel: "Reset password",
		ButtonURL:   resetURL,
		LinkHint:    "Or open:",
		LinkURL:     resetURL,
		Footer:      "For security, this link expires in 7 minutes. Never share this email with anyone.",
	})
	if err != nil {
		logger.L().Error().Err(err).
			Str("mail_driver", "smtp").
			Str("mail_outcome", "failed").
			Str("kind", "password_reset").
			Str("to", toEmail).
			Msg("mail: smtp email HTML render failed")
		return err
	}
	return s.send(toEmail, subject, plain, bodyHTML, "password_reset")
}

func (s *SMTPMailer) validate() error {
	if strings.TrimSpace(s.cfg.MailFromAddress) == "" {
		return fmt.Errorf("mail: MAIL_FROM_ADDRESS is required when MAIL_MAILER=smtp")
	}
	if strings.TrimSpace(s.cfg.MailHost) == "" {
		return fmt.Errorf("mail: MAIL_HOST is required when MAIL_MAILER=smtp")
	}
	return nil
}

func (s *SMTPMailer) send(to, subject, plainText, htmlBody string, kind string) error {
	from := formatMailFrom(s.cfg.MailFromName, s.cfg.MailFromAddress)
	fromAddr := strings.TrimSpace(s.cfg.MailFromAddress)

	port := strings.TrimSpace(s.cfg.MailPort)
	if port == "" {
		port = "587"
	}
	portNum, err := strconv.Atoi(port)
	if err != nil || portNum < 1 || portNum > 65535 {
		portErr := fmt.Errorf("mail: invalid MAIL_PORT %q", s.cfg.MailPort)
		logger.L().Error().Err(portErr).
			Str("mail_driver", "smtp").
			Str("mail_outcome", "failed").
			Str("kind", kind).
			Str("to", to).
			Msg("mail: smtp invalid port")
		return portErr
	}

	host := strings.TrimSpace(s.cfg.MailHost)
	addr := net.JoinHostPort(host, strconv.Itoa(portNum))
	auth := smtp.PlainAuth("", strings.TrimSpace(s.cfg.MailUsername), s.cfg.MailPassword, host)
	if s.cfg.MailUsername == "" || s.cfg.MailPassword == "" {
		auth = nil
	}
	raw := buildMultipartMessage(from, to, subject, plainText, htmlBody)
	if err := smtp.SendMail(addr, auth, fromAddr, []string{to}, raw); err != nil {
		logger.L().Error().Err(err).
			Str("mail_driver", "smtp").
			Str("mail_outcome", "failed").
			Str("kind", kind).
			Str("to", to).
			Str("subject", subject).
			Str("from", fromAddr).
			Str("smtp_host", host).
			Msg("mail: smtp send failed")
		return err
	}
	logger.L().Info().
		Str("mail_driver", "smtp").
		Str("mail_outcome", "success").
		Str("kind", kind).
		Str("to", to).
		Str("subject", subject).
		Str("from", fromAddr).
		Str("smtp_host", host).
		Msg("mail: transactional email sent successfully")
	return nil
}

func formatMailFrom(name, address string) string {
	name = strings.TrimSpace(name)
	address = strings.TrimSpace(address)
	if name == "" {
		return address
	}
	return fmt.Sprintf("%s <%s>", name, address)
}

func randomBoundary() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "b_mix_alternative_fallback"
	}
	return "b_" + hex.EncodeToString(b)
}

func buildMultipartMessage(from, to, subject, plainText, htmlBody string) []byte {
	boundary := randomBoundary()
	var sb strings.Builder
	sb.WriteString("From: ")
	sb.WriteString(from)
	sb.WriteString("\r\nTo: ")
	sb.WriteString(to)
	sb.WriteString("\r\nSubject: ")
	sb.WriteString(subject)
	sb.WriteString("\r\nMIME-Version: 1.0\r\n")
	sb.WriteString("Content-Type: multipart/alternative; boundary=")
	sb.WriteString(boundary)
	sb.WriteString("\r\n\r\n")

	sb.WriteString("--")
	sb.WriteString(boundary)
	sb.WriteString("\r\nContent-Type: text/plain; charset=UTF-8\r\nContent-Transfer-Encoding: 8bit\r\n\r\n")
	sb.WriteString(plainText)
	if !strings.HasSuffix(plainText, "\n") {
		sb.WriteString("\r\n")
	}

	sb.WriteString("\r\n--")
	sb.WriteString(boundary)
	sb.WriteString("\r\nContent-Type: text/html; charset=UTF-8\r\nContent-Transfer-Encoding: 8bit\r\n\r\n")
	sb.WriteString(htmlBody)
	sb.WriteString("\r\n\r\n--")
	sb.WriteString(boundary)
	sb.WriteString("--\r\n")

	return []byte(sb.String())
}
