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

// PublicAuthCSRFRequired protects public auth endpoints using synchronizer token check.
func PublicAuthCSRFRequired(store *auth.PublicAuthCSRFStore, headerName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.L()

		headerToken := strings.TrimSpace(c.GetHeader(headerName))
		if headerToken == "" {
			log.Warn().
				Str("path", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Str("header_name", headerName).
				Msg("public auth csrf failed: missing csrf header")
			response.FailureWithAbort(c, http.StatusForbidden, "csrf token is required", "csrf token is required")
			return
		}

		if err := store.ValidateAndConsume(headerToken); err != nil {
			if errors.Is(err, auth.ErrPublicAuthCSRFExpired) {
				log.Warn().
					Str("path", c.Request.URL.Path).
					Str("method", c.Request.Method).
					Msg("public auth csrf failed: token expired")
				response.FailureWithAbort(c, http.StatusForbidden, "csrf token expired", "csrf token expired")
				return
			}

			log.Warn().
				Err(err).
				Str("path", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Msg("public auth csrf failed: token invalid")
			response.FailureWithAbort(c, http.StatusForbidden, "csrf validation failed", "csrf validation failed")
			return
		}

		c.Next()
	}
}
