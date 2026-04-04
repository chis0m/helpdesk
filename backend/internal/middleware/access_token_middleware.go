package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"helpdesk/backend/internal/auth"
	"helpdesk/backend/internal/logger"
	"helpdesk/backend/internal/response"
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

		// Only access tokens should pass protected-route middleware.
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
