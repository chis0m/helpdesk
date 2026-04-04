package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"helpdesk/backend/internal/auth"
	"helpdesk/backend/internal/config"
	"helpdesk/backend/internal/logger"
	"helpdesk/backend/internal/models"
	"helpdesk/backend/internal/repositories"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrInvalidRefreshToken = errors.New("invalid refresh token")
var ErrInvalidSession = errors.New("invalid session")

type LoginResult struct {
	User   *models.User
	Tokens auth.TokenPair
	CSRF   auth.CSRFToken
}

type AuthService struct {
	cfg         config.Config
	tokenMaker  auth.MakerInterface
	userRepo    *repositories.UserRepository
	sessionRepo *repositories.AuthSessionRepository
}

func NewAuthService(
	cfg config.Config,
	tokenMaker auth.MakerInterface,
	userRepo *repositories.UserRepository,
	sessionRepo *repositories.AuthSessionRepository,
) *AuthService {
	return &AuthService{
		cfg:         cfg,
		tokenMaker:  tokenMaker,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (s *AuthService) Login(email, password string) (*LoginResult, error) {
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

	if _, err := s.sessionRepo.Create(user.UUID, sessionID, refreshPayload.Jti, refreshPayload.Exp); err != nil {
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
