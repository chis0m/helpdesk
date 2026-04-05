package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"helpdesk/backend/internal/logger"
	"helpdesk/backend/internal/repositories"
	"helpdesk/backend/internal/response"
)

// CSRFRequired enforces CSRF checks for state-changing requests only.
func CSRFRequired(sessionRepo *repositories.AuthSessionRepository, headerName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if isSafeMethod(c.Request.Method) {
			c.Next()
			return
		}

		log := logger.L()

		sessionIDRaw, ok := c.Get(CtxSessionID)
		if !ok {
			log.Warn().
				Str("path", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Msg("csrf validation failed: missing session id in context")
			response.FailureWithAbort(c, http.StatusForbidden, "csrf validation failed", "csrf validation failed")
			return
		}

		sessionID, ok := sessionIDRaw.(string)
		if !ok || sessionID == "" {
			log.Warn().
				Str("path", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Msg("csrf validation failed: invalid session id in context")
			response.FailureWithAbort(c, http.StatusForbidden, "csrf validation failed", "csrf validation failed")
			return
		}

		sessionUUID, err := uuid.Parse(sessionID)
		if err != nil {
			log.Warn().
				Err(err).
				Str("path", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Msg("csrf validation failed: malformed session id")
			response.FailureWithAbort(c, http.StatusForbidden, "csrf validation failed", "csrf validation failed")
			return
		}

		session, err := sessionRepo.GetActiveBySessionID(sessionUUID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Warn().
					Err(err).
					Str("session_id", sessionID).
					Msg("csrf validation failed: active session not found")
				response.FailureWithAbort(c, http.StatusForbidden, "csrf validation failed", "csrf validation failed")
				return
			}
			log.Error().
				Err(err).
				Str("session_id", sessionID).
				Msg("csrf validation failed: session lookup error")
			response.FailureWithAbort(c, http.StatusForbidden, "csrf validation failed", "csrf validation failed")
			return
		}

		token := strings.TrimSpace(c.GetHeader(headerName))
		if token == "" {
			log.Warn().
				Str("path", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Str("header_name", headerName).
				Msg("csrf validation failed: missing csrf header")
			response.FailureWithAbort(c, http.StatusForbidden, "csrf token is required", "csrf token is required")
			return
		}

		if session.CSRFToken == nil || session.CSRFExpiresAt == nil {
			log.Warn().
				Str("path", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Str("session_id", sessionID).
				Msg("csrf validation failed: csrf token not initialized on session")
			response.FailureWithAbort(c, http.StatusForbidden, "csrf token is not initialized", "csrf token is not initialized")
			return
		}

		if session.CSRFExpiresAt.UTC().Before(time.Now().UTC()) {
			log.Warn().
				Str("path", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Str("session_id", sessionID).
				Msg("csrf validation failed: csrf token expired")
			response.FailureWithAbort(c, http.StatusForbidden, "csrf token expired", "csrf token expired")
			return
		}

		// VULN-05: Broken CSRF (session token not verified) — header not compared to session.CSRFToken; any non-empty value passes.
		c.Next()
	}
}

func isSafeMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	default:
		return false
	}
}
