package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"helpdesk/backend/internal/audit"
	"helpdesk/backend/internal/logger"
	"helpdesk/backend/internal/middleware"
	"helpdesk/backend/internal/models"
	"helpdesk/backend/internal/repositories"
	"helpdesk/backend/internal/requests"
	"helpdesk/backend/internal/response"
)

type AuditLogController struct {
	auditLogRepo *repositories.AuditLogRepository
}

func NewAuditLogController(auditLogRepo *repositories.AuditLogRepository) *AuditLogController {
	return &AuditLogController{auditLogRepo: auditLogRepo}
}

// SEC-05: Read-only audit log listing for admin / super_admin (append-only store i.e only insert operations are allowed, no update or delete operations are allowed).
func (a *AuditLogController) List(c *gin.Context) {
	log := logger.L()

	roleValue, ok := c.Get(middleware.CtxUserRole)
	if !ok {
		log.Warn().Msg("list audit logs failed: missing user role in context")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}
	roleStr, ok := roleValue.(string)
	if !ok || roleStr == "" {
		log.Warn().Msg("list audit logs failed: invalid role in context")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}
	actorRole := models.UserRole(roleStr)
	if actorRole != models.RoleAdmin && actorRole != models.RoleSuperAdmin {
		log.Warn().Str("role", roleStr).Msg("list audit logs failed: admin or super_admin access required")
		response.FailureWithAbort(c, http.StatusForbidden, "admin or super_admin access required", "admin or super_admin access required")
		return
	}

	var query requests.ListAuditLogsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		log.Warn().Err(err).Msg("list audit logs failed: invalid query params")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid query parameters", "invalid query parameters")
		return
	}

	page := query.Page
	if page < 1 {
		page = 1
	}
	limit := query.Limit
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	rows, total, err := a.auditLogRepo.List(c.Request.Context(), page, limit)
	if err != nil {
		log.Error().Err(err).Msg("list audit logs failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	items := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		var meta interface{}
		if len(row.Metadata) > 0 {
			if err := json.Unmarshal(row.Metadata, &meta); err != nil {
				meta = string(row.Metadata)
			}
		}
		item := gin.H{
			"id":              row.ID,
			"created_at":      row.CreatedAt,
			"http_method":     row.HTTPMethod,
			"path":            row.Path,
			"action":          row.Action,
			"success":         row.Success,
			"resource_type":   row.ResourceType,
			"resource_id":     row.ResourceID,
			"metadata":        meta,
			"actor_user_uuid": row.ActorUserUUID,
			"session_id":      row.SessionID,
			"ip":              row.IP,
			"user_agent":      row.UserAgent,
		}
		if row.ErrorCode != nil {
			item["error_code"] = *row.ErrorCode
		}
		items = append(items, item)
	}

	// SEC-05: Record that an authorized principal exported a page of the append-only audit trail.
	audit.Write(c, a.auditLogRepo, audit.Event{
		Action:       audit.ActionAdminAuditLogList,
		Success:      true,
		ResourceType: audit.Str(audit.ResourceTypeAuditLog),
		Metadata: map[string]interface{}{
			"page":      page,
			"limit":     limit,
			"row_count": len(rows),
			"total":     total,
		},
	})

	response.Success(c, http.StatusOK, gin.H{
		"items": items,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	}, "audit logs")
}
