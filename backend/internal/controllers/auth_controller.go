package controllers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"helpdesk/backend/internal/config"
	"helpdesk/backend/internal/logger"
	"helpdesk/backend/internal/requests"
	"helpdesk/backend/internal/response"
	"helpdesk/backend/internal/services"
)

type AuthController struct {
	cfg         config.Config
	authService *services.AuthService
}

func NewAuthController(cfg config.Config, authService *services.AuthService) *AuthController {
	return &AuthController{
		cfg:         cfg,
		authService: authService,
	}
}

func (a *AuthController) Login(c *gin.Context) {
	log := logger.L()

	var req requests.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid request payload", "invalid request payload")
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	result, err := a.authService.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			log.Warn().
				Str("email", strings.TrimSpace(req.Email)).
				Msg("login failed: invalid credentials")
			response.FailureWithAbort(c, http.StatusUnauthorized, "invalid email or password", "invalid email or password")
			return
		}
		log.Error().Err(err).Msg("login failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	setAuthCookies(c, a.cfg, result.AccessToken, result.RefreshToken, result.AccessExpires, result.RefreshExpires)

	response.Success(c, http.StatusOK, gin.H{
		"user_uuid":             result.User.UUID.String(),
		"email":                 result.User.Email,
		"role":                  result.User.Role,
		"must_change_password":  result.User.MustChangePassword,
		"access_expires_at_utc": result.AccessExpires.UTC(),
	}, "login successful")
}
