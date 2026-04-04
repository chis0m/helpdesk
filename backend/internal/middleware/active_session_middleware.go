package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"helpdesk/backend/internal/logger"
	"helpdesk/backend/internal/repositories"
	"helpdesk/backend/internal/response"
)

// ActiveSessionRequired is an optional second-layer check:
// after token verification, it ensures session_id still exists
// and is not revoked in auth_sessions.
func ActiveSessionRequired(sessionRepo *repositories.AuthSessionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.L()

		sessionIDRaw, ok := c.Get(CtxSessionID)
		if !ok {
			log.Warn().
				Str("path", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Msg("auth failed: missing session id in context")
			response.FailureWithAbort(c, http.StatusUnauthorized, "invalid session", "invalid session")
			return
		}

		sessionID, ok := sessionIDRaw.(string)
		if !ok || sessionID == "" {
			log.Warn().
				Str("path", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Msg("auth failed: invalid session id type in context")
			response.FailureWithAbort(c, http.StatusUnauthorized, "invalid session", "invalid session")
			return
		}

		sessionUUID, err := uuid.Parse(sessionID)
		if err != nil {
			log.Warn().
				Err(err).
				Str("path", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Msg("auth failed: malformed session id")
			response.FailureWithAbort(c, http.StatusUnauthorized, "invalid session", "invalid session")
			return
		}

		if _, err := sessionRepo.GetActiveBySessionID(sessionUUID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Warn().
					Err(err).
					Str("path", c.Request.URL.Path).
					Str("method", c.Request.Method).
					Msg("auth failed: inactive or revoked session")
				response.FailureWithAbort(c, http.StatusUnauthorized, "session is not active", "session is not active")
				return
			}
			log.Error().
				Err(err).
				Str("path", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Msg("auth failed: session lookup error")
			response.FailureWithAbort(c, http.StatusUnauthorized, "session is not active", "session is not active")
			return
		}

		c.Next()
	}
}
