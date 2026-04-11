package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"helpdesk/backend/internal/audit"
	"helpdesk/backend/internal/config"
	"helpdesk/backend/internal/logger"
	"helpdesk/backend/internal/middleware"
	"helpdesk/backend/internal/models"
	"helpdesk/backend/internal/repositories"
	"helpdesk/backend/internal/requests"
	"helpdesk/backend/internal/response"
	"helpdesk/backend/internal/services"
)

type InviteController struct {
	cfg           config.Config
	inviteService *services.InviteService
	userService   *services.UserService
	auditLogRepo  *repositories.AuditLogRepository
}

func NewInviteController(cfg config.Config, inviteService *services.InviteService, userService *services.UserService, auditLogRepo *repositories.AuditLogRepository) *InviteController {
	return &InviteController{
		cfg:           cfg,
		inviteService: inviteService,
		userService:   userService,
		auditLogRepo:  auditLogRepo,
	}
}

func (ic *InviteController) CreateStaffInvite(c *gin.Context) {
	log := logger.L()

	roleValue, ok := c.Get(middleware.CtxUserRole)
	if !ok {
		log.Warn().Msg("create staff invite failed: missing user role in context")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}
	roleStr, ok := roleValue.(string)
	actorRole := models.UserRole(roleStr)
	if !ok || (actorRole != models.RoleAdmin && actorRole != models.RoleSuperAdmin) {
		log.Warn().Str("role", roleStr).Msg("create staff invite failed: admin or super_admin access required")
		response.FailureWithAbort(c, http.StatusForbidden, "admin or super_admin access required", "admin or super_admin access required")
		return
	}

	userUUIDRaw, ok := c.Get(middleware.CtxUserUUID)
	if !ok {
		log.Warn().Msg("create staff invite failed: missing user uuid in context")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}
	userUUIDStr, ok := userUUIDRaw.(string)
	if !ok || userUUIDStr == "" {
		log.Warn().Msg("create staff invite failed: invalid user uuid in context")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}

	actor, err := ic.userService.GetByUUIDString(userUUIDStr)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Err(err).Msg("create staff invite failed: actor not found")
			response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
			return
		}
		log.Warn().Err(err).Msg("create staff invite failed: invalid actor identity")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}

	var req requests.CreateStaffInviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("create staff invite failed: invalid request payload")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid request payload", "invalid request payload")
		return
	}

	inv, inviteURL, err := ic.inviteService.CreateStaffInvite(actor.ID, actorRole, req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInviteForbidden):
			log.Warn().Msg("create staff invite failed: forbidden")
			response.FailureWithAbort(c, http.StatusForbidden, "forbidden", "forbidden")
		case errors.Is(err, services.ErrInviteAdminForbidden):
			log.Warn().Str("actor_role", string(actorRole)).Msg("create staff invite failed: only admin or super_admin may invite admin")
			response.FailureWithAbort(c, http.StatusForbidden, "only admin or super_admin may invite users with role admin", "only admin or super_admin may invite users with role admin")
		case errors.Is(err, services.ErrInviteEmailTaken):
			log.Warn().Str("email", req.Email).Msg("create staff invite failed: email already registered")
			response.FailureWithAbort(c, http.StatusConflict, "email already registered", "email already registered")
		case errors.Is(err, services.ErrInvitePendingExists):
			log.Warn().Str("email", req.Email).Msg("create staff invite failed: pending invite exists")
			response.FailureWithAbort(c, http.StatusConflict, "pending invite already exists for this email", "pending invite already exists for this email")
		default:
			log.Error().Err(err).Msg("create staff invite failed")
			response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		}
		return
	}

	delivery := "log"
	if ic.cfg.UseSMTPMail() {
		delivery = "smtp"
	}
	iid := inv.ID
	audit.Write(c, ic.auditLogRepo, audit.Event{
		Action:       audit.ActionInviteStaffCreate,
		Success:      true,
		ActorUserID:  &actor.ID,
		ResourceType: audit.Str(audit.ResourceTypeInvite),
		ResourceID:   &iid,
		Metadata: map[string]interface{}{
			"target_email": inv.Email,
			"target_role":  string(inv.TargetRole),
			"delivery":     delivery,
		},
	})

	payload := gin.H{
		"invite_id":      inv.ID,
		"email":          inv.Email,
		"expires_at_utc": inv.ExpiresAt.UTC(),
		"target_role":    inv.TargetRole,
	}
	if ic.cfg.UseSMTPMail() {
		payload["delivery"] = "smtp"
		payload["notice"] = "Invite email was handed to the configured SMTP server. Check the recipient inbox (or Mailtrap sandbox)."
	} else {
		payload["delivery"] = "log"
		payload["invite_url"] = inviteURL
		payload["notice"] = "MAIL_MAILER is not smtp — no real email was sent. Copy invite_url for the staff member, or set MAIL_MAILER=smtp and MAIL_* for SMTP delivery."
	}

	response.Success(c, http.StatusCreated, payload, "staff invite created")
}

func (ic *InviteController) VerifyInvite(c *gin.Context) {
	log := logger.L()

	token := c.Query("token")
	result, err := ic.inviteService.VerifyInvite(token)
	if err != nil {
		log.Error().Err(err).Msg("verify invite failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	if !result.Valid {
		response.Success(c, http.StatusOK, gin.H{"valid": false}, "invite verification")
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"valid":      true,
		"email":      result.Email,
		"first_name": result.FirstName,
		"last_name":  result.LastName,
	}, "invite verification")
}

func (ic *InviteController) AcceptInvite(c *gin.Context) {
	log := logger.L()

	var req requests.AcceptInviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("accept invite failed: invalid request payload")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid request payload", "invalid request payload")
		return
	}

	user, err := ic.inviteService.AcceptInvite(req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInviteNotFound):
			log.Warn().Msg("accept invite failed: invalid token")
			response.FailureWithAbort(c, http.StatusBadRequest, "invalid or expired invite", "invalid or expired invite")
		case errors.Is(err, services.ErrInviteExpired):
			log.Warn().Msg("accept invite failed: expired")
			response.FailureWithAbort(c, http.StatusBadRequest, "invalid or expired invite", "invalid or expired invite")
		case errors.Is(err, services.ErrInviteUsed):
			log.Warn().Msg("accept invite failed: already used")
			response.FailureWithAbort(c, http.StatusBadRequest, "invite already used", "invite already used")
		case errors.Is(err, services.ErrInviteEmailTaken):
			log.Warn().Msg("accept invite failed: email taken")
			response.FailureWithAbort(c, http.StatusConflict, "email already registered", "email already registered")
		default:
			log.Error().Err(err).Msg("accept invite failed")
			response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		}
		return
	}

	uid := user.ID
	audit.Write(c, ic.auditLogRepo, audit.Event{
		Action:       audit.ActionInviteAccepted,
		Success:      true,
		ActorUserID:  &uid,
		ResourceType: audit.Str(audit.ResourceTypeUser),
		ResourceID:   &uid,
		Metadata: map[string]interface{}{
			"email":  user.Email,
			"role":   user.Role,
			"method": "invite",
		},
	})

	response.Success(c, http.StatusCreated, gin.H{
		"user_id":     user.ID,
		"user_uuid":   user.UUID.String(),
		"email":       user.Email,
		"role":        user.Role,
		"redirect_to": "/login",
	}, "invite accepted; account created")
}
