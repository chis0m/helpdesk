package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"helpdesk/backend/internal/auth"
	"helpdesk/backend/internal/logger"
	"helpdesk/backend/internal/response"
)

func RefreshTokenRequired(maker auth.MakerInterface, refreshCookieName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.L()

		token, err := c.Cookie(refreshCookieName)
		if err != nil {
			log.Warn().
				Err(err).
				Str("path", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Str("cookie_name", refreshCookieName).
				Msg("refresh auth failed: missing refresh cookie")
			response.FailureWithAbort(c, http.StatusUnauthorized, "refresh token required", "refresh token required")
			return
		}

		payload, err := maker.VerifyToken(token)
		if err != nil {
			log.Warn().
				Err(err).
				Str("path", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Msg("refresh auth failed: token verification failed")
			response.FailureWithAbort(c, http.StatusUnauthorized, "invalid refresh token", "invalid refresh token")
			return
		}

		if payload.Type != auth.TokenTypeRefresh {
			log.Warn().
				Str("path", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Str("token_type", string(payload.Type)).
				Msg("refresh auth failed: invalid token type")
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
