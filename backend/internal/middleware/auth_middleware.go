package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gin-gonic/gin"

	"helpdesk/backend/internal/auth"
	"helpdesk/backend/internal/logger"
	"helpdesk/backend/internal/repositories"
	"helpdesk/backend/internal/response"
)

const (
	CtxAuthPayload = "auth_payload"
	CtxUserUUID    = "user_uuid"
	CtxUserRole    = "user_role"
	CtxSessionID   = "session_id"
	CtxTokenJTI    = "token_jti"
)

func AuthRequired(maker auth.MakerInterface, accessCookieName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.L()

		token, err := c.Cookie(accessCookieName)
		if err != nil {
			log.Warn().
				Err(err).
				Str("path", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Str("cookie_name", accessCookieName).
				Msg("auth failed: missing access cookie")
			response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
			return
		}

		payload, err := maker.VerifyToken(token)
		if err != nil {
			log.Warn().
				Err(err).
				Str("path", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Msg("auth failed: token verification failed")
			response.FailureWithAbort(c, http.StatusUnauthorized, "invalid access token", "invalid access token")
			return
		}

		// Only token of type "access" should pass protected-route middleware.
		// Refresh tokens are for token renewal endpoints only.
		if payload.Type != auth.TokenTypeAccess {
			log.Warn().
				Str("path", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Str("token_type", string(payload.Type)).
				Msg("auth failed: invalid token type")
			response.FailureWithAbort(c, http.StatusUnauthorized, "invalid token type", "invalid token type")
			return
		}

		c.Set(CtxAuthPayload, payload)
		c.Set(CtxUserUUID, payload.Sub)
		c.Set(CtxUserRole, payload.Role)
		c.Set(CtxSessionID, payload.SessID)
		c.Set(CtxTokenJTI, payload.Jti)

		c.Next()
	}
}

// ActiveSessionRequired
// after token verification, it ensures session_id still exists and not revoked.
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
			log.Warn().
				Err(err).
				Str("path", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Msg("auth failed: inactive or revoked session")
			response.FailureWithAbort(c, http.StatusUnauthorized, "session is not active", "session is not active")
			return
		}

		c.Next()
	}
}
