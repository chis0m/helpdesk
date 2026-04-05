package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"helpdesk/backend/internal/auth"
	"helpdesk/backend/internal/config"
	"helpdesk/backend/internal/logger"
	"helpdesk/backend/internal/mail"
	"helpdesk/backend/internal/models"
	"helpdesk/backend/internal/repositories"
	"helpdesk/backend/internal/requests"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrInvalidRefreshToken = errors.New("invalid refresh token")
var ErrInvalidSession = errors.New("invalid session")
var ErrSignupFailed = errors.New("unable to complete signup")
var ErrInvalidPassword = errors.New("invalid password")
var ErrPasswordResetInvalid = errors.New("invalid password reset token")
var ErrPasswordResetExpired = errors.New("password reset token expired")
var ErrPasswordResetUsed = errors.New("password reset token already used")
var ErrSessionRevokeNotFound = errors.New("session not found or access denied")

type LoginResult struct {
	User   *models.User
	Tokens auth.TokenPair
	CSRF   auth.CSRFToken
}

type AuthService struct {
	cfg               config.Config
	tokenMaker        auth.MakerInterface
	userRepo          *repositories.UserRepository
	sessionRepo       *repositories.AuthSessionRepository
	passwordResetRepo *repositories.PasswordResetRepository
	resetNotifier     mail.PasswordResetNotifier
}

func NewAuthService(
	cfg config.Config,
	tokenMaker auth.MakerInterface,
	userRepo *repositories.UserRepository,
	sessionRepo *repositories.AuthSessionRepository,
	passwordResetRepo *repositories.PasswordResetRepository,
	resetNotifier mail.PasswordResetNotifier,
) *AuthService {
	return &AuthService{
		cfg:               cfg,
		tokenMaker:        tokenMaker,
		userRepo:          userRepo,
		sessionRepo:       sessionRepo,
		passwordResetRepo: passwordResetRepo,
		resetNotifier:     resetNotifier,
	}
}

// AuthSessionListItem is a safe subset of auth_sessions for the current user.
type AuthSessionListItem struct {
	SessionID string
	CreatedAt time.Time
	UserAgent *string
	IP        *string
	IsCurrent bool
}

func hashPasswordResetToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func generatePasswordResetRawToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func truncateSessionUserAgent(ua string) *string {
	ua = strings.TrimSpace(ua)
	if ua == "" {
		return nil
	}
	const maxRunes = 512
	runes := []rune(ua)
	if len(runes) > maxRunes {
		ua = string(runes[:maxRunes])
	}
	return &ua
}

func truncateSessionIP(ip string) *string {
	ip = strings.TrimSpace(ip)
	if ip == "" {
		return nil
	}
	const maxLen = 45
	if len(ip) > maxLen {
		ip = ip[:maxLen]
	}
	return &ip
}

func (s *AuthService) Login(email, password string, userAgent string, clientIP string) (*LoginResult, error) {
	log := logger.L()

	email = strings.ToLower(strings.TrimSpace(email))
	password = strings.TrimSpace(password)

	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	ok, err := auth.VerifyPassword(password, user.PasswordHash)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrInvalidCredentials
	}

	sessionID := uuid.New()

	accessClaims := auth.Claims{
		Issuer:    s.cfg.TokenIssuer(),
		Subject:   user.UUID.String(),
		Role:      string(user.Role),
		Audience:  s.cfg.TokenAudience(),
		Duration:  s.cfg.AccessTokenTTL(),
		TokenType: auth.TokenTypeAccess,
		SessionID: sessionID.String(),
	}
	accessToken, accessPayload, err := s.tokenMaker.CreateAccessToken(accessClaims)
	if err != nil {
		return nil, fmt.Errorf("create access token: %w", err)
	}

	refreshClaims := auth.Claims{
		Issuer:    s.cfg.TokenIssuer(),
		Subject:   user.UUID.String(),
		Role:      string(user.Role),
		Audience:  s.cfg.TokenAudience(),
		Duration:  s.cfg.RefreshTokenTTL(),
		TokenType: auth.TokenTypeRefresh,
		SessionID: sessionID.String(),
	}
	refreshToken, refreshPayload, err := s.tokenMaker.CreateRefreshToken(refreshClaims)
	if err != nil {
		return nil, fmt.Errorf("create refresh token: %w", err)
	}

	if _, err := s.sessionRepo.Create(
		user.UUID,
		sessionID,
		refreshPayload.Jti,
		refreshPayload.Exp,
		truncateSessionUserAgent(userAgent),
		truncateSessionIP(clientIP),
	); err != nil {
		return nil, fmt.Errorf("create auth session: %w", err)
	}

	csrfToken, err := s.IssueCSRFToken(sessionID)
	if err != nil {
		return nil, fmt.Errorf("issue csrf token: %w", err)
	}

	log.Info().
		Str("user_uuid", user.UUID.String()).
		Str("session_id", sessionID.String()).
		Msg("user login successful")

	return &LoginResult{
		User: user,
		Tokens: auth.TokenPair{
			AccessToken:    accessToken,
			RefreshToken:   refreshToken,
			AccessExpires:  accessPayload.Exp,
			RefreshExpires: refreshPayload.Exp,
			CSRFToken:      csrfToken.Token,
			CSRFExpiresAt:  csrfToken.ExpiresAt,
		},
		CSRF: *csrfToken,
	}, nil
}

func (s *AuthService) Signup(input requests.SignupRequest) (*models.User, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))
	firstName := strings.TrimSpace(input.FirstName)
	lastName := strings.TrimSpace(input.LastName)
	password := strings.TrimSpace(input.Password)

	var middleName *string
	if input.MiddleName != nil {
		trimmed := strings.TrimSpace(*input.MiddleName)
		if trimmed != "" {
			middleName = &trimmed
		}
	}

	if _, err := s.userRepo.GetByEmail(email); err == nil {
		return nil, ErrSignupFailed
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	isActive := true
	mustChangePassword := false
	changedAt := time.Now().UTC()
	createInput := requests.CreateUserInput{
		Email:              email,
		PasswordHash:       passwordHash,
		FirstName:          firstName,
		LastName:           lastName,
		MiddleName:         middleName,
		Role:               models.RoleUser,
		IsActive:           &isActive,
		MustChangePassword: &mustChangePassword,
		PasswordChangedAt:  &changedAt,
	}

	user, err := s.userRepo.Create(createInput)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Refresh(refreshPayload *auth.Payload) (*auth.TokenPair, error) {
	log := logger.L()

	if refreshPayload == nil || refreshPayload.Type != auth.TokenTypeRefresh {
		return nil, ErrInvalidRefreshToken
	}

	sessionID, err := uuid.Parse(refreshPayload.SessID)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	session, err := s.sessionRepo.GetActiveBySessionID(sessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidRefreshToken
		}
		return nil, err
	}

	if session.RefreshExpiresAt.UTC().Before(time.Now().UTC()) {
		return nil, ErrInvalidRefreshToken
	}

	if session.RefreshJTI != refreshPayload.Jti {
		_ = s.sessionRepo.RevokeBySessionID(sessionID)
		log.Warn().
			Str("session_id", sessionID.String()).
			Msg("refresh token replay detected; session revoked")
		return nil, ErrInvalidRefreshToken
	}

	accessClaims := auth.Claims{
		Issuer:    s.cfg.TokenIssuer(),
		Subject:   refreshPayload.Sub,
		Role:      refreshPayload.Role,
		Audience:  s.cfg.TokenAudience(),
		Duration:  s.cfg.AccessTokenTTL(),
		TokenType: auth.TokenTypeAccess,
		SessionID: sessionID.String(),
	}
	accessToken, accessPayload, err := s.tokenMaker.CreateAccessToken(accessClaims)
	if err != nil {
		return nil, fmt.Errorf("create access token: %w", err)
	}

	refreshClaims := auth.Claims{
		Issuer:    s.cfg.TokenIssuer(),
		Subject:   refreshPayload.Sub,
		Role:      refreshPayload.Role,
		Audience:  s.cfg.TokenAudience(),
		Duration:  s.cfg.RefreshTokenTTL(),
		TokenType: auth.TokenTypeRefresh,
		SessionID: sessionID.String(),
	}
	refreshToken, newRefreshPayload, err := s.tokenMaker.CreateRefreshToken(refreshClaims)
	if err != nil {
		return nil, fmt.Errorf("create refresh token: %w", err)
	}

	if err := s.sessionRepo.RotateRefreshJTI(sessionID, newRefreshPayload.Jti, newRefreshPayload.Exp); err != nil {
		return nil, fmt.Errorf("rotate refresh token: %w", err)
	}

	csrfToken, err := s.GetOrIssueCSRFToken(session)
	if err != nil {
		return nil, fmt.Errorf("get or issue csrf token: %w", err)
	}

	return &auth.TokenPair{
		AccessToken:    accessToken,
		RefreshToken:   refreshToken,
		AccessExpires:  accessPayload.Exp,
		RefreshExpires: newRefreshPayload.Exp,
		CSRFToken:      csrfToken.Token,
		CSRFExpiresAt:  csrfToken.ExpiresAt,
	}, nil
}

func (s *AuthService) IssueCSRFToken(sessionID uuid.UUID) (*auth.CSRFToken, error) {
	tokenID := uuid.NewString()
	expiresAt := time.Now().UTC().Add(s.cfg.CSRFTTL())

	if err := s.sessionRepo.UpsertCSRFToken(sessionID, tokenID, expiresAt); err != nil {
		return nil, err
	}

	return &auth.CSRFToken{
		Token:     tokenID,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *AuthService) IssueCSRFTokenBySessionID(sessionID string) (*auth.CSRFToken, error) {
	sessionUUID, err := uuid.Parse(strings.TrimSpace(sessionID))
	if err != nil {
		return nil, ErrInvalidSession
	}

	session, err := s.sessionRepo.GetActiveBySessionID(sessionUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidSession
		}
		return nil, err
	}

	return s.GetOrIssueCSRFToken(session)
}

func (s *AuthService) GetOrIssueCSRFToken(session *models.AuthSession) (*auth.CSRFToken, error) {
	if session == nil {
		return nil, ErrInvalidSession
	}

	if session.CSRFToken != nil && session.CSRFExpiresAt != nil && session.CSRFExpiresAt.UTC().After(time.Now().UTC()) {
		return &auth.CSRFToken{
			Token:     *session.CSRFToken,
			ExpiresAt: session.CSRFExpiresAt.UTC(),
		}, nil
	}

	return s.IssueCSRFToken(session.SessionID)
}

func (s *AuthService) Logout(sessionID string) error {
	sessionUUID, err := uuid.Parse(strings.TrimSpace(sessionID))
	if err != nil {
		return ErrInvalidSession
	}
	if err := s.sessionRepo.RevokeBySessionID(sessionUUID); err != nil {
		return err
	}
	return nil
}

func (s *AuthService) GetMe(userUUID string) (*models.User, error) {
	parsed, err := uuid.Parse(strings.TrimSpace(userUUID))
	if err != nil {
		return nil, ErrInvalidSession
	}
	return s.userRepo.GetByUUID(parsed)
}

func (s *AuthService) ChangePassword(userUUID, currentPassword, newPassword string) error {
	parsed, err := uuid.Parse(strings.TrimSpace(userUUID))
	if err != nil {
		return ErrInvalidSession
	}

	user, err := s.userRepo.GetByUUID(parsed)
	if err != nil {
		return err
	}

	ok, err := auth.VerifyPassword(strings.TrimSpace(currentPassword), user.PasswordHash)
	if err != nil {
		return err
	}
	if !ok {
		return ErrInvalidPassword
	}

	newHash, err := auth.HashPassword(strings.TrimSpace(newPassword))
	if err != nil {
		return err
	}

	if err := s.userRepo.UpdatePasswordByUUID(parsed, newHash, false, time.Now().UTC()); err != nil {
		return err
	}
	return nil
}

// RequestPasswordReset always succeeds from the caller's perspective when the input is valid.
// If the email is registered, a row is stored and the reset URL is logged (CA: no SMTP).
func (s *AuthService) RequestPasswordReset(email string) error {
	email = strings.ToLower(strings.TrimSpace(email))
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	now := time.Now().UTC()
	if err := s.passwordResetRepo.InvalidateUnusedForUser(user.UUID, now); err != nil {
		return err
	}

	rawToken, err := generatePasswordResetRawToken()
	if err != nil {
		return err
	}

	row := &models.PasswordReset{
		UserUUID:  user.UUID,
		TokenHash: hashPasswordResetToken(rawToken),
		ExpiresAt: now.Add(s.cfg.PasswordResetTTL()),
	}
	if err := s.passwordResetRepo.Create(row); err != nil {
		return err
	}

	base := strings.TrimSuffix(strings.TrimSpace(s.cfg.FrontendURL), "/")
	resetURL := base + "/reset-password?token=" + url.QueryEscape(rawToken)
	return s.resetNotifier.SendPasswordReset(user.Email, resetURL)
}

func (s *AuthService) CompletePasswordReset(rawToken, newPassword string) error {
	rawToken = strings.TrimSpace(rawToken)
	if len(rawToken) != 64 {
		return ErrPasswordResetInvalid
	}
	if _, err := hex.DecodeString(rawToken); err != nil {
		return ErrPasswordResetInvalid
	}

	row, err := s.passwordResetRepo.GetByTokenHash(hashPasswordResetToken(rawToken))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPasswordResetInvalid
		}
		return err
	}

	now := time.Now().UTC()
	if row.UsedAt != nil {
		return ErrPasswordResetUsed
	}
	if !row.ExpiresAt.UTC().After(now) {
		return ErrPasswordResetExpired
	}

	newHash, err := auth.HashPassword(strings.TrimSpace(newPassword))
	if err != nil {
		return err
	}

	if err := s.userRepo.UpdatePasswordByUUID(row.UserUUID, newHash, false, now); err != nil {
		return err
	}
	if err := s.passwordResetRepo.MarkUsed(row.ID, now); err != nil {
		return err
	}
	if err := s.sessionRepo.RevokeAllActiveForUser(row.UserUUID); err != nil {
		return err
	}
	return nil
}

func (s *AuthService) ListSessionsForUser(userUUIDStr, currentSessionIDStr string) ([]AuthSessionListItem, error) {
	userUUID, err := uuid.Parse(strings.TrimSpace(userUUIDStr))
	if err != nil {
		return nil, ErrInvalidSession
	}
	currentSessionID, err := uuid.Parse(strings.TrimSpace(currentSessionIDStr))
	if err != nil {
		return nil, ErrInvalidSession
	}

	rows, err := s.sessionRepo.ListActiveByUserUUID(userUUID)
	if err != nil {
		return nil, err
	}

	out := make([]AuthSessionListItem, 0, len(rows))
	for _, sess := range rows {
		out = append(out, AuthSessionListItem{
			SessionID: sess.SessionID.String(),
			CreatedAt: sess.CreatedAt.UTC(),
			UserAgent: sess.UserAgent,
			IP:        sess.IP,
			IsCurrent: sess.SessionID == currentSessionID,
		})
	}
	return out, nil
}

// RevokeSessionForUser revokes a session that belongs to the user. Returns whether the revoked session was the current one (caller may clear cookies).
func (s *AuthService) RevokeSessionForUser(userUUIDStr, currentSessionIDStr, targetSessionIDStr string) (isCurrent bool, err error) {
	userUUID, err := uuid.Parse(strings.TrimSpace(userUUIDStr))
	if err != nil {
		return false, ErrInvalidSession
	}
	currentSessionID, err := uuid.Parse(strings.TrimSpace(currentSessionIDStr))
	if err != nil {
		return false, ErrInvalidSession
	}
	targetSessionID, err := uuid.Parse(strings.TrimSpace(targetSessionIDStr))
	if err != nil {
		return false, ErrInvalidSession
	}

	revoked, err := s.sessionRepo.RevokeActiveSessionForUser(targetSessionID, userUUID)
	if err != nil {
		return false, err
	}
	if !revoked {
		return false, ErrSessionRevokeNotFound
	}
	return targetSessionID == currentSessionID, nil
}

func (s *AuthService) RevokeOtherSessionsForUser(userUUIDStr, currentSessionIDStr string) error {
	userUUID, err := uuid.Parse(strings.TrimSpace(userUUIDStr))
	if err != nil {
		return ErrInvalidSession
	}
	currentSessionID, err := uuid.Parse(strings.TrimSpace(currentSessionIDStr))
	if err != nil {
		return ErrInvalidSession
	}
	return s.sessionRepo.RevokeOtherActiveSessionsForUser(currentSessionID, userUUID)
}
