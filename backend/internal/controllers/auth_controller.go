package controllers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"helpdesk/backend/internal/auth"
	"helpdesk/backend/internal/config"
	"helpdesk/backend/internal/logger"
	"helpdesk/backend/internal/middleware"
	"helpdesk/backend/internal/requests"
	"helpdesk/backend/internal/response"
	"helpdesk/backend/internal/services"
)

type AuthController struct {
	cfg            config.Config
	authService    *services.AuthService
	publicAuthCSRF *auth.PublicAuthCSRFStore
}

func NewAuthController(
	cfg config.Config,
	authService *services.AuthService,
	publicAuthCSRF *auth.PublicAuthCSRFStore,
) *AuthController {
	return &AuthController{
		cfg:            cfg,
		authService:    authService,
		publicAuthCSRF: publicAuthCSRF,
	}
}

func (a *AuthController) Login(c *gin.Context) {
	log := logger.L()

	var req requests.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("login failed: invalid request payload")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid request payload", "invalid request payload")
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	result, err := a.authService.Login(req.Email, req.Password, c.Request.UserAgent(), c.ClientIP())
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

	// VULN-01: Weak session cookie flags — HttpOnly/Secure false, SameSite=None.
	setAuthCookies(c, a.cfg, result.Tokens)

	response.Success(c, http.StatusOK, gin.H{
		"user_uuid":             result.User.UUID.String(),
		"email":                 result.User.Email,
		"role":                  result.User.Role,
		"must_change_password":  result.User.MustChangePassword,
		"access_expires_at_utc": result.Tokens.AccessExpires.UTC(),
		"csrf_token":            result.CSRF.Token,
		"csrf_expires_at_utc":   result.CSRF.ExpiresAt.UTC(),
	}, "login successful")
}

func (a *AuthController) PublicAuthCSRFToken(c *gin.Context) {
	result := a.publicAuthCSRF.Issue()

	response.Success(c, http.StatusOK, gin.H{
		"csrf_token":          result.Token,
		"csrf_expires_at_utc": result.ExpiresAt.UTC(),
	}, "public auth csrf token ready")
}

func (a *AuthController) Signup(c *gin.Context) {
	log := logger.L()

	var req requests.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("signup failed: invalid request payload")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid request payload", "invalid request payload")
		return
	}

	user, err := a.authService.Signup(req)
	if err != nil {
		if errors.Is(err, services.ErrSignupFailed) {
			log.Warn().
				Str("email", strings.TrimSpace(req.Email)).
				Msg("signup failed")
			response.FailureWithAbort(c, http.StatusBadRequest, "unable to complete signup", "unable to complete signup")
			return
		}
		log.Error().Err(err).Msg("signup failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	response.Success(c, http.StatusCreated, gin.H{
		"user_uuid":   user.UUID.String(),
		"email":       user.Email,
		"redirect_to": "/login",
	}, "signup successful")
}

func (a *AuthController) Refresh(c *gin.Context) {
	log := logger.L()

	value, ok := c.Get(middleware.CtxAuthPayload)
	if !ok {
		log.Warn().Msg("refresh failed: missing auth payload in context")
		response.FailureWithAbort(c, http.StatusUnauthorized, "invalid refresh token", "invalid refresh token")
		return
	}

	refreshPayload, ok := value.(*auth.Payload)
	if !ok || refreshPayload == nil {
		log.Warn().Msg("refresh failed: invalid auth payload type in context")
		response.FailureWithAbort(c, http.StatusUnauthorized, "invalid refresh token", "invalid refresh token")
		return
	}

	result, err := a.authService.Refresh(refreshPayload)
	if err != nil {
		if errors.Is(err, services.ErrInvalidRefreshToken) {
			log.Warn().
				Str("session_id", refreshPayload.SessID).
				Msg("refresh failed: invalid refresh token")
			response.FailureWithAbort(c, http.StatusUnauthorized, "invalid refresh token", "invalid refresh token")
			return
		}
		log.Error().Err(err).Msg("refresh failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	// VULN-01: Weak session cookie flags — HttpOnly/Secure false, SameSite=None.
	setAuthCookies(c, a.cfg, *result)
	response.Success(c, http.StatusOK, gin.H{
		"access_expires_at_utc": result.AccessExpires.UTC(),
		"csrf_token":            result.CSRFToken,
		"csrf_expires_at_utc":   result.CSRFExpiresAt.UTC(),
	}, "refresh successful")
}

func (a *AuthController) CSRFToken(c *gin.Context) {
	log := logger.L()

	sessionIDRaw, ok := c.Get(middleware.CtxSessionID)
	if !ok {
		log.Warn().Msg("csrf token issue failed: missing session id in context")
		response.FailureWithAbort(c, http.StatusUnauthorized, "invalid session", "invalid session")
		return
	}

	sessionID, ok := sessionIDRaw.(string)
	if !ok || strings.TrimSpace(sessionID) == "" {
		log.Warn().Msg("csrf token issue failed: invalid session id in context")
		response.FailureWithAbort(c, http.StatusUnauthorized, "invalid session", "invalid session")
		return
	}

	result, err := a.authService.IssueCSRFTokenBySessionID(sessionID)
	if err != nil {
		if errors.Is(err, services.ErrInvalidSession) {
			log.Warn().
				Str("session_id", strings.TrimSpace(sessionID)).
				Msg("csrf token issue failed: invalid session")
			response.FailureWithAbort(c, http.StatusUnauthorized, "invalid session", "invalid session")
			return
		}
		log.Error().Err(err).Msg("csrf token issue failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	log.Info().
		Str("session_id", strings.TrimSpace(sessionID)).
		Msg("csrf token issued")
	response.Success(c, http.StatusOK, gin.H{
		"csrf_token":          result.Token,
		"csrf_expires_at_utc": result.ExpiresAt.UTC(),
	}, "csrf token ready")
}

func (a *AuthController) Logout(c *gin.Context) {
	log := logger.L()

	sessionIDRaw, ok := c.Get(middleware.CtxSessionID)
	if !ok {
		log.Warn().Msg("logout failed: missing session id in context")
		response.FailureWithAbort(c, http.StatusUnauthorized, "invalid session", "invalid session")
		return
	}

	sessionID, ok := sessionIDRaw.(string)
	if !ok || strings.TrimSpace(sessionID) == "" {
		log.Warn().Msg("logout failed: invalid session id in context")
		response.FailureWithAbort(c, http.StatusUnauthorized, "invalid session", "invalid session")
		return
	}

	if err := a.authService.Logout(sessionID); err != nil {
		if errors.Is(err, services.ErrInvalidSession) {
			log.Warn().Str("session_id", strings.TrimSpace(sessionID)).Msg("logout failed: invalid session")
			response.FailureWithAbort(c, http.StatusUnauthorized, "invalid session", "invalid session")
			return
		}
		log.Error().Err(err).Str("session_id", strings.TrimSpace(sessionID)).Msg("logout failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	clearAuthCookies(c)
	response.Success(c, http.StatusOK, gin.H{
		"redirect_to": "/login",
	}, "logout successful")
}

func (a *AuthController) Me(c *gin.Context) {
	log := logger.L()

	userUUIDRaw, ok := c.Get(middleware.CtxUserUUID)
	if !ok {
		log.Warn().Msg("me failed: missing user uuid in context")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}

	userUUID, ok := userUUIDRaw.(string)
	if !ok || strings.TrimSpace(userUUID) == "" {
		log.Warn().Msg("me failed: invalid user uuid in context")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}

	user, err := a.authService.GetMe(userUUID)
	if err != nil {
		if errors.Is(err, services.ErrInvalidSession) {
			log.Warn().Str("user_uuid", strings.TrimSpace(userUUID)).Msg("me failed: invalid user uuid")
			response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Str("user_uuid", strings.TrimSpace(userUUID)).Msg("me failed: user not found")
			response.FailureWithAbort(c, http.StatusNotFound, "user not found", "user not found")
			return
		}
		log.Error().Err(err).Str("user_uuid", strings.TrimSpace(userUUID)).Msg("me failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"user_id":              user.ID,
		"user_uuid":            user.UUID.String(),
		"email":                user.Email,
		"first_name":           user.FirstName,
		"last_name":            user.LastName,
		"middle_name":          user.MiddleName,
		"role":                 user.Role,
		"is_active":            user.IsActive,
		"must_change_password": user.MustChangePassword,
	}, "me fetched successfully")
}

func (a *AuthController) ChangePassword(c *gin.Context) {
	log := logger.L()

	userUUIDRaw, ok := c.Get(middleware.CtxUserUUID)
	if !ok {
		log.Warn().Msg("change password failed: missing user uuid in context")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}

	userUUID, ok := userUUIDRaw.(string)
	if !ok || strings.TrimSpace(userUUID) == "" {
		log.Warn().Msg("change password failed: invalid user uuid in context")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}

	var req requests.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("change password failed: invalid request payload")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid request payload", "invalid request payload")
		return
	}

	if err := a.authService.ChangePassword(userUUID, req.CurrentPassword, req.NewPassword); err != nil {
		if errors.Is(err, services.ErrInvalidSession) {
			log.Warn().Str("user_uuid", strings.TrimSpace(userUUID)).Msg("change password failed: invalid session")
			response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
			return
		}
		if errors.Is(err, services.ErrInvalidPassword) {
			log.Warn().Str("user_uuid", strings.TrimSpace(userUUID)).Msg("change password failed: invalid current password")
			response.FailureWithAbort(c, http.StatusBadRequest, "current password is invalid", "current password is invalid")
			return
		}
		log.Error().Err(err).Str("user_uuid", strings.TrimSpace(userUUID)).Msg("change password failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	response.Success(c, http.StatusOK, nil, "password changed successfully")
}

func (a *AuthController) ForgotPassword(c *gin.Context) {
	log := logger.L()

	var req requests.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("forgot password failed: invalid request payload")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid request payload", "invalid request payload")
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	if err := a.authService.RequestPasswordReset(req.Email); err != nil {
		log.Error().Err(err).Str("email", req.Email).Msg("forgot password failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	// Same response whether or not the email exists (avoid account enumeration).
	response.Success(c, http.StatusOK, nil,
		"If that email is registered, a password reset link was written to the server logs (CA: no email delivery).")
}

func (a *AuthController) ResetPassword(c *gin.Context) {
	log := logger.L()

	var req requests.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("reset password failed: invalid request payload")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid request payload", "invalid request payload")
		return
	}

	if err := a.authService.CompletePasswordReset(req.Token, req.NewPassword); err != nil {
		if errors.Is(err, services.ErrPasswordResetInvalid) ||
			errors.Is(err, services.ErrPasswordResetExpired) ||
			errors.Is(err, services.ErrPasswordResetUsed) {
			log.Warn().Err(err).Msg("reset password failed: invalid or expired token")
			response.FailureWithAbort(c, http.StatusBadRequest, "invalid or expired reset token", "invalid or expired reset token")
			return
		}
		log.Error().Err(err).Msg("reset password failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"redirect_to": "/login"}, "password reset successful")
}

func (a *AuthController) ListSessions(c *gin.Context) {
	log := logger.L()

	userUUIDRaw, ok := c.Get(middleware.CtxUserUUID)
	if !ok {
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}
	userUUID, ok := userUUIDRaw.(string)
	if !ok || strings.TrimSpace(userUUID) == "" {
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}

	sessionIDRaw, ok := c.Get(middleware.CtxSessionID)
	if !ok {
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}
	sessionID, ok := sessionIDRaw.(string)
	if !ok || strings.TrimSpace(sessionID) == "" {
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}

	items, err := a.authService.ListSessionsForUser(userUUID, sessionID)
	if err != nil {
		if errors.Is(err, services.ErrInvalidSession) {
			response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
			return
		}
		log.Error().Err(err).Msg("list sessions failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	out := make([]gin.H, 0, len(items))
	for _, it := range items {
		var ua any
		if it.UserAgent != nil {
			ua = *it.UserAgent
		}
		var ip any
		if it.IP != nil {
			ip = *it.IP
		}
		out = append(out, gin.H{
			"session_id": it.SessionID,
			"created_at": it.CreatedAt,
			"user_agent": ua,
			"ip":         ip,
			"is_current": it.IsCurrent,
		})
	}

	response.Success(c, http.StatusOK, gin.H{"items": out}, "sessions fetched")
}

func (a *AuthController) RevokeSession(c *gin.Context) {
	log := logger.L()

	userUUIDRaw, ok := c.Get(middleware.CtxUserUUID)
	if !ok {
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}
	userUUID, ok := userUUIDRaw.(string)
	if !ok || strings.TrimSpace(userUUID) == "" {
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}

	sessionIDRaw, ok := c.Get(middleware.CtxSessionID)
	if !ok {
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}
	currentSessionID, ok := sessionIDRaw.(string)
	if !ok || strings.TrimSpace(currentSessionID) == "" {
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}

	target := strings.TrimSpace(c.Param("session_id"))
	if _, err := uuid.Parse(target); err != nil {
		log.Warn().Err(err).Str("session_id", target).Msg("revoke session failed: invalid session id")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid session id", "invalid session id")
		return
	}

	isCurrent, err := a.authService.RevokeSessionForUser(userUUID, currentSessionID, target)
	if err != nil {
		if errors.Is(err, services.ErrInvalidSession) {
			response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
			return
		}
		if errors.Is(err, services.ErrSessionRevokeNotFound) {
			log.Warn().Str("target_session_id", target).Msg("revoke session failed: not found or denied")
			response.FailureWithAbort(c, http.StatusNotFound, "session not found", "session not found")
			return
		}
		log.Error().Err(err).Msg("revoke session failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	if isCurrent {
		clearAuthCookies(c)
	}

	response.Success(c, http.StatusOK, gin.H{
		"revoked_session_id": target,
		"logged_out":         isCurrent,
	}, "session revoked")
}

func (a *AuthController) RevokeMyOtherSessions(c *gin.Context) {
	log := logger.L()

	userUUIDRaw, ok := c.Get(middleware.CtxUserUUID)
	if !ok {
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}
	userUUID, ok := userUUIDRaw.(string)
	if !ok || strings.TrimSpace(userUUID) == "" {
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}

	sessionIDRaw, ok := c.Get(middleware.CtxSessionID)
	if !ok {
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}
	currentSessionID, ok := sessionIDRaw.(string)
	if !ok || strings.TrimSpace(currentSessionID) == "" {
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}

	if err := a.authService.RevokeOtherSessionsForUser(userUUID, currentSessionID); err != nil {
		if errors.Is(err, services.ErrInvalidSession) {
			response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
			return
		}
		log.Error().Err(err).Msg("revoke other sessions failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	response.Success(c, http.StatusOK, nil, "other sessions revoked")
}
