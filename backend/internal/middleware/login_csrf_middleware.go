package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"helpdesk/backend/internal/auth"
	"helpdesk/backend/internal/logger"
	"helpdesk/backend/internal/response"
)

// LoginCSRFRequired protects login endpoint using synchronizer token check.
func LoginCSRFRequired(store *auth.LoginCSRFStore, headerName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.L()

		headerToken := strings.TrimSpace(c.GetHeader(headerName))
		if headerToken == "" {
			log.Warn().
				Str("path", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Str("header_name", headerName).
				Msg("login csrf failed: missing csrf header")
			response.FailureWithAbort(c, http.StatusForbidden, "csrf token is required", "csrf token is required")
			return
		}

		if err := store.ValidateAndConsume(headerToken); err != nil {
			if errors.Is(err, auth.ErrLoginCSRFExpired) {
				log.Warn().
					Str("path", c.Request.URL.Path).
					Str("method", c.Request.Method).
					Msg("login csrf failed: token expired")
				response.FailureWithAbort(c, http.StatusForbidden, "csrf token expired", "csrf token expired")
				return
			}

			log.Warn().
				Err(err).
				Str("path", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Msg("login csrf failed: token invalid")
			response.FailureWithAbort(c, http.StatusForbidden, "csrf validation failed", "csrf validation failed")
			return
		}

		c.Next()
	}
}
