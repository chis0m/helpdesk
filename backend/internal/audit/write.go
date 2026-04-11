package audit

import (
	"context"
	"encoding/json"

	"github.com/gin-gonic/gin"

	"helpdesk/backend/internal/logger"
	"helpdesk/backend/internal/middleware"
	"helpdesk/backend/internal/models"
)

// LogAppender persists audit rows (implemented by repositories.AuditLogRepository).
type LogAppender interface {
	Create(ctx context.Context, row *models.AuditLog) error
}

// Write persists an append-only audit row. Failures are logged and do not affect the HTTP response.
func Write(c *gin.Context, repo LogAppender, e Event) {
	EnrichFromGin(c, &e)
	row, err := eventToModel(&e)
	if err != nil {
		logger.L().Error().Err(err).Str("action", e.Action).Msg("audit: failed to build audit row")
		return
	}
	if err := repo.Create(c.Request.Context(), row); err != nil {
		logger.L().Error().Err(err).Str("action", e.Action).Msg("audit: insert failed")
	}
}

// EnrichFromGin fills HTTP and auth context. It does not overwrite ActorUserID if already set.
func EnrichFromGin(c *gin.Context, e *Event) {
	e.HTTPMethod = c.Request.Method
	e.Path = c.Request.URL.Path

	if v, ok := c.Get(middleware.CtxSessionID); ok {
		if s, ok := v.(string); ok && s != "" {
			e.SessionID = &s
		}
	}
	if v, ok := c.Get(middleware.CtxTokenJTI); ok {
		if s, ok := v.(string); ok && s != "" {
			e.TokenJTI = &s
		}
	}

	ip := c.ClientIP()
	e.IP = &ip
	if ua := c.GetHeader("User-Agent"); ua != "" {
		e.UserAgent = &ua
	}
}

func eventToModel(e *Event) (*models.AuditLog, error) {
	var meta []byte
	if len(e.Metadata) > 0 {
		b, err := json.Marshal(e.Metadata)
		if err != nil {
			return nil, err
		}
		meta = b
	}

	row := &models.AuditLog{
		ActorUserID: e.ActorUserID,
		SessionID:   e.SessionID,
		TokenJTI:    e.TokenJTI,
		HTTPMethod:  e.HTTPMethod,
		Path:        e.Path,
		Action:      e.Action,
		Success:       e.Success,
		ResourceType:  e.ResourceType,
		ResourceID:    e.ResourceID,
		IP:            e.IP,
		UserAgent:     e.UserAgent,
		Metadata:      meta,
	}
	if e.ErrorCode != "" {
		ec := e.ErrorCode
		row.ErrorCode = &ec
	}
	return row, nil
}

// Str returns a pointer to s (for ResourceType and similar).
func Str(s string) *string { return &s }
