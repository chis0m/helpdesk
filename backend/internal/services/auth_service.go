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

type LoginResult struct {
	User           *models.User
	AccessToken    string
	RefreshToken   string
	AccessExpires  time.Time
	RefreshExpires time.Time
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

	log.Info().
		Str("user_uuid", user.UUID.String()).
		Str("session_id", sessionID.String()).
		Msg("user login successful")

	return &LoginResult{
		User:           user,
		AccessToken:    accessToken,
		RefreshToken:   refreshToken,
		AccessExpires:  accessPayload.Exp,
		RefreshExpires: refreshPayload.Exp,
	}, nil
}
