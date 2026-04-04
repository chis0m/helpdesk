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

type LoginResult struct {
	User   *models.User
	Tokens auth.TokenPair
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
		log.Error().Err(err).Msg("get user by email failed")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	ok, err := auth.VerifyPassword(password, user.PasswordHash)
	if err != nil {
		log.Error().Err(err).Msg("verify password failed")
		return nil, err
	}
	if !ok {
		log.Error().Err(ErrInvalidCredentials).Msg("invalid credentials")
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
		log.Error().Err(err).Msg("create access token failed")
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
		log.Error().Err(err).Msg("create refresh token failed")
		return nil, fmt.Errorf("create refresh token: %w", err)
	}

	if _, err := s.sessionRepo.Create(user.UUID, sessionID, refreshPayload.Jti, refreshPayload.Exp); err != nil {
		log.Error().Err(err).Msg("create auth session failed")
		return nil, fmt.Errorf("create auth session: %w", err)
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
		},
	}, nil
}

func (s *AuthService) Refresh(refreshPayload *auth.Payload) (*auth.TokenPair, error) {
	log := logger.L()

	if refreshPayload == nil || refreshPayload.Type != auth.TokenTypeRefresh {
		log.Error().Err(ErrInvalidRefreshToken).Msg("refresh token required")
		return nil, ErrInvalidRefreshToken
	}

	sessionID, err := uuid.Parse(refreshPayload.SessID)
	if err != nil {
		log.Error().Err(err).Msg("parse session ID failed")
		return nil, ErrInvalidRefreshToken
	}

	session, err := s.sessionRepo.GetActiveBySessionID(sessionID)
	if err != nil {
		log.Error().Err(err).Msg("get active session by session ID failed")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidRefreshToken
		}
		return nil, err
	}

	if session.RefreshExpiresAt.UTC().Before(time.Now().UTC()) {
		log.Error().Err(ErrInvalidRefreshToken).Msg("refresh token expired")
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

	return &auth.TokenPair{
		AccessToken:    accessToken,
		RefreshToken:   refreshToken,
		AccessExpires:  accessPayload.Exp,
		RefreshExpires: newRefreshPayload.Exp,
	}, nil
}
