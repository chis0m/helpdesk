package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"helpdesk/backend/internal/logger"
	"helpdesk/backend/internal/middleware"
	"helpdesk/backend/internal/models"
	"helpdesk/backend/internal/repositories"
	"helpdesk/backend/internal/requests"
	"helpdesk/backend/internal/response"
	"helpdesk/backend/internal/services"
)

type TicketController struct {
	ticketService *services.TicketService
	userRepo      *repositories.UserRepository
}

func NewTicketController(ticketService *services.TicketService, userRepo *repositories.UserRepository) *TicketController {
	return &TicketController{ticketService: ticketService, userRepo: userRepo}
}

// VULN-03: Weak input validation / stored XSS risk — Create binds ticket JSON with no HTML/script sanitization.
func (t *TicketController) Create(c *gin.Context) {
	log := logger.L()

	userUUID, ok := c.Get(middleware.CtxUserUUID)
	if !ok {
		log.Warn().Msg("create ticket failed: missing authenticated user")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}
	userUUIDStr, ok := userUUID.(string)
	if !ok || userUUIDStr == "" {
		log.Warn().Msg("create ticket failed: invalid authenticated user value")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}

	var req requests.CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("create ticket failed: invalid request payload")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid request payload", "invalid request payload")
		return
	}

	ticket, err := t.ticketService.CreateByUserUUID(userUUIDStr, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Err(err).Msg("create ticket failed: reporter not found")
			response.FailureWithAbort(c, http.StatusNotFound, "user not found", "user not found")
			return
		}
		log.Error().Err(err).Msg("create ticket failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	users, err := t.loadTicketUsers([]models.Ticket{*ticket})
	if err != nil {
		log.Error().Err(err).Msg("create ticket failed: load user names")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}
	response.Success(c, http.StatusCreated, formatTicket(ticket, users), "ticket created")
}

func (t *TicketController) List(c *gin.Context) {
	log := logger.L()
	// Non-admin roles are scoped to tickets they report or are assigned to; admins may list all.

	var query requests.ListTicketsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		log.Warn().Err(err).Msg("list tickets failed: invalid query params")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid query parameters", "invalid query parameters")
		return
	}

	filter := requests.ListTicketsFilter{
		Page:           query.Page,
		Limit:          query.Limit,
		Status:         query.Status,
		Category:       query.Category,
		ReporterUserID: query.ReporterUserID,
		AssignedUserID: query.AssignedUserID,
	}

	userUUID, role, ok := getAuthenticatedUser(c)
	if !ok {
		log.Warn().Msg("list tickets failed: missing authenticated user context")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}

	tickets, total, err := t.ticketService.ListForActor(userUUID, role, filter)
	if err != nil {
		if errors.Is(err, services.ErrTicketListForbidden) {
			log.Warn().Str("user_uuid", userUUID).Msg("list tickets failed: forbidden filters for role")
			response.FailureWithAbort(c, http.StatusForbidden, "forbidden", "forbidden")
			return
		}
		log.Error().Err(err).Msg("list tickets failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	users, err := t.loadTicketUsers(tickets)
	if err != nil {
		log.Error().Err(err).Msg("list tickets failed: load user names")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	result := make([]gin.H, 0, len(tickets))
	for _, ticket := range tickets {
		tk := ticket
		result = append(result, formatTicket(&tk, users))
	}

	response.Success(c, http.StatusOK, gin.H{
		"items": result,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	}, "tickets fetched")
}

// VULN-06: SQL injection (ticket keyword search) — forwards q to repository Raw SQL without parameter binding.
func (t *TicketController) Search(c *gin.Context) {
	log := logger.L()

	q := strings.TrimSpace(c.Query("q"))
	if q == "" {
		response.FailureWithAbort(c, http.StatusBadRequest, "query parameter q is required", "query parameter q is required")
		return
	}

	tickets, err := t.ticketService.SearchTicketsUnsafe(q)
	if err != nil {
		log.Error().Err(err).Str("q", q).Msg("ticket search failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	users, err := t.loadTicketUsers(tickets)
	if err != nil {
		log.Error().Err(err).Msg("ticket search failed: load user names")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	result := make([]gin.H, 0, len(tickets))
	for _, ticket := range tickets {
		tk := ticket
		result = append(result, formatTicket(&tk, users))
	}

	response.Success(c, http.StatusOK, gin.H{
		"items": result,
		"query": q,
	}, "tickets search completed")
}

func (t *TicketController) GetByID(c *gin.Context) {
	log := logger.L()
	// VULN-02: IDOR on tickets and comments — no reporter/assignee/admin check on ticket id.

	ticketID, ok := parseUintID(c.Param("id"))
	if !ok {
		log.Warn().Str("ticket_id", c.Param("id")).Msg("get ticket failed: invalid id")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid id", "invalid id")
		return
	}

	ticket, err := t.ticketService.GetByID(ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Uint64("ticket_id", ticketID).Msg("get ticket failed: ticket not found")
			response.FailureWithAbort(c, http.StatusNotFound, "ticket not found", "ticket not found")
			return
		}
		log.Error().Err(err).Uint64("ticket_id", ticketID).Msg("get ticket failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	users, err := t.loadTicketUsers([]models.Ticket{*ticket})
	if err != nil {
		log.Error().Err(err).Uint64("ticket_id", ticketID).Msg("get ticket failed: load user names")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}
	response.Success(c, http.StatusOK, formatTicket(ticket, users), "ticket fetched")
}

func (t *TicketController) UpdateByID(c *gin.Context) {
	log := logger.L()
	// VULN-02: IDOR on tickets and comments — no reporter/assignee/admin check on ticket id.

	ticketID, ok := parseUintID(c.Param("id"))
	if !ok {
		log.Warn().Str("ticket_id", c.Param("id")).Msg("update ticket failed: invalid id")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid id", "invalid id")
		return
	}

	var req requests.UpdateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Uint64("ticket_id", ticketID).Msg("update ticket failed: invalid request payload")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid request payload", "invalid request payload")
		return
	}

	ticket, err := t.ticketService.UpdateByID(ticketID, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Uint64("ticket_id", ticketID).Msg("update ticket failed: ticket not found")
			response.FailureWithAbort(c, http.StatusNotFound, "ticket not found", "ticket not found")
			return
		}
		log.Error().Err(err).Uint64("ticket_id", ticketID).Msg("update ticket failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	users, err := t.loadTicketUsers([]models.Ticket{*ticket})
	if err != nil {
		log.Error().Err(err).Uint64("ticket_id", ticketID).Msg("update ticket failed: load user names")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}
	response.Success(c, http.StatusOK, formatTicket(ticket, users), "ticket updated")
}

func (t *TicketController) UpdateStatus(c *gin.Context) {
	log := logger.L()
	// VULN-02: IDOR on tickets and comments — no reporter/assignee/admin check on ticket id.

	ticketID, ok := parseUintID(c.Param("id"))
	if !ok {
		log.Warn().Str("ticket_id", c.Param("id")).Msg("update ticket status failed: invalid id")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid id", "invalid id")
		return
	}

	var req requests.UpdateTicketStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Uint64("ticket_id", ticketID).Msg("update ticket status failed: invalid request payload")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid request payload", "invalid request payload")
		return
	}

	ticket, err := t.ticketService.UpdateStatus(ticketID, req.Status)
	if err != nil {
		if errors.Is(err, services.ErrInvalidTicketStatusTransition) {
			log.Warn().Uint64("ticket_id", ticketID).Str("requested_status", string(req.Status)).Msg("update ticket status failed: invalid status transition")
			response.FailureWithAbort(c, http.StatusBadRequest, "invalid status transition", "invalid status transition")
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Uint64("ticket_id", ticketID).Msg("update ticket status failed: ticket not found")
			response.FailureWithAbort(c, http.StatusNotFound, "ticket not found", "ticket not found")
			return
		}
		log.Error().Err(err).Uint64("ticket_id", ticketID).Msg("update ticket status failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	users, err := t.loadTicketUsers([]models.Ticket{*ticket})
	if err != nil {
		log.Error().Err(err).Uint64("ticket_id", ticketID).Msg("update ticket status failed: load user names")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}
	response.Success(c, http.StatusOK, formatTicket(ticket, users), "ticket status updated")
}

func (t *TicketController) Assign(c *gin.Context) {
	log := logger.L()

	_, actorRole, authOk := getAuthenticatedUser(c)
	if !authOk {
		log.Warn().Msg("assign ticket failed: unauthenticated")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}

	ticketID, ok := parseUintID(c.Param("id"))
	if !ok {
		log.Warn().Str("ticket_id", c.Param("id")).Msg("assign ticket failed: invalid id")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid id", "invalid id")
		return
	}

	var req requests.AssignTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Uint64("ticket_id", ticketID).Msg("assign ticket failed: invalid request payload")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid request payload", "invalid request payload")
		return
	}

	if req.AssignedUserID == nil {
		response.FailureWithAbort(c, http.StatusBadRequest, "assigned_user_id required", "assigned_user_id required")
		return
	}

	ticket, err := t.ticketService.Assign(actorRole, ticketID, req.AssignedUserID, req.Unassign)
	if err != nil {
		if errors.Is(err, services.ErrTicketAssignForbidden) {
			response.FailureWithAbort(c, http.StatusForbidden, "forbidden", "admin or super_admin required")
			return
		}
		if errors.Is(err, services.ErrTicketAssigneeInvalidRole) {
			response.FailureWithAbort(c, http.StatusBadRequest, "invalid assignee", "assignee must be staff or admin")
			return
		}
		if errors.Is(err, services.ErrTicketAlreadyAssignedToUser) {
			response.FailureWithAbort(c, http.StatusBadRequest, "already assigned", "ticket is already assigned to this user")
			return
		}
		if errors.Is(err, services.ErrTicketAlreadyUnassigned) {
			response.FailureWithAbort(c, http.StatusBadRequest, "not assigned", "ticket has no assignee")
			return
		}
		if errors.Is(err, services.ErrTicketAssignMissingAssigneeID) {
			response.FailureWithAbort(c, http.StatusBadRequest, "assigned_user_id required", "assigned_user_id required")
			return
		}
		if errors.Is(err, services.ErrTicketUnassignAssigneeMismatch) {
			response.FailureWithAbort(c, http.StatusBadRequest, "assignee mismatch", "selected user is not the current assignee")
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Uint64("ticket_id", ticketID).Msg("assign ticket failed: ticket or user not found")
			response.FailureWithAbort(c, http.StatusNotFound, "resource not found", "resource not found")
			return
		}
		log.Error().Err(err).Uint64("ticket_id", ticketID).Msg("assign ticket failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	users, err := t.loadTicketUsers([]models.Ticket{*ticket})
	if err != nil {
		log.Error().Err(err).Uint64("ticket_id", ticketID).Msg("assign ticket failed: load user names")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}
	response.Success(c, http.StatusOK, formatTicket(ticket, users), "ticket assignment updated")
}

func (t *TicketController) DeleteByID(c *gin.Context) {
	log := logger.L()
	// VULN-02: IDOR on tickets and comments — no reporter/assignee/admin check on ticket id.

	ticketID, ok := parseUintID(c.Param("id"))
	if !ok {
		log.Warn().Str("ticket_id", c.Param("id")).Msg("delete ticket failed: invalid id")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid id", "invalid id")
		return
	}

	if _, err := t.ticketService.GetByID(ticketID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Uint64("ticket_id", ticketID).Msg("delete ticket failed: ticket not found")
			response.FailureWithAbort(c, http.StatusNotFound, "ticket not found", "ticket not found")
			return
		}
		log.Error().Err(err).Uint64("ticket_id", ticketID).Msg("delete ticket failed: pre-delete lookup failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	if err := t.ticketService.DeleteByID(ticketID); err != nil {
		log.Error().Err(err).Uint64("ticket_id", ticketID).Msg("delete ticket failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"ticket_id": ticketID}, "ticket deleted")
}

func (t *TicketController) AddComment(c *gin.Context) {
	log := logger.L()
	// VULN-02: IDOR on tickets and comments — no reporter/assignee/admin check on ticket id.

	ticketID, ok := parseUintID(c.Param("id"))
	if !ok {
		log.Warn().Str("ticket_id", c.Param("id")).Msg("add ticket comment failed: invalid ticket id")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid id", "invalid id")
		return
	}

	userUUID, _, ok := getAuthenticatedUser(c)
	if !ok {
		log.Warn().Uint64("ticket_id", ticketID).Msg("add ticket comment failed: missing authenticated user context")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}

	var req requests.CreateTicketCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Uint64("ticket_id", ticketID).Msg("add ticket comment failed: invalid request payload")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid request payload", "invalid request payload")
		return
	}

	comment, err := t.ticketService.AddComment(ticketID, userUUID, req.Body)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Uint64("ticket_id", ticketID).Msg("add ticket comment failed: ticket or actor not found")
			response.FailureWithAbort(c, http.StatusNotFound, "resource not found", "resource not found")
			return
		}
		log.Error().Err(err).Uint64("ticket_id", ticketID).Msg("add ticket comment failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	response.Success(c, http.StatusCreated, gin.H{
		"comment_id":     comment.ID,
		"ticket_id":      comment.TicketID,
		"author_user_id": comment.AuthorUserID,
		"body":           comment.Body,
		"created_at":     comment.CreatedAt,
		"updated_at":     comment.UpdatedAt,
	}, "ticket comment created")
}

func (t *TicketController) ListComments(c *gin.Context) {
	log := logger.L()
	// VULN-02: IDOR on tickets and comments — no reporter/assignee/admin check on ticket id.

	ticketID, ok := parseUintID(c.Param("id"))
	if !ok {
		log.Warn().Str("ticket_id", c.Param("id")).Msg("list ticket comments failed: invalid ticket id")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid id", "invalid id")
		return
	}

	comments, err := t.ticketService.ListComments(ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Uint64("ticket_id", ticketID).Msg("list ticket comments failed: ticket not found")
			response.FailureWithAbort(c, http.StatusNotFound, "ticket not found", "ticket not found")
			return
		}
		log.Error().Err(err).Uint64("ticket_id", ticketID).Msg("list ticket comments failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	// Nil slice JSON-marshals as null; clients expect data.items to always be an array.
	if comments == nil {
		comments = []models.TicketCommentWithAuthor{}
	}

	response.Success(c, http.StatusOK, gin.H{"items": comments}, "ticket comments fetched")
}

func (t *TicketController) UpdateComment(c *gin.Context) {
	log := logger.L()
	// VULN-02: IDOR on tickets and comments — ticket id in path not access-controlled; comment edit is author/admin only.

	ticketID, ok := parseUintID(c.Param("id"))
	if !ok {
		log.Warn().Str("ticket_id", c.Param("id")).Msg("update ticket comment failed: invalid ticket id")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid id", "invalid id")
		return
	}
	commentID, ok := parseUintID(c.Param("commentId"))
	if !ok {
		log.Warn().Str("comment_id", c.Param("commentId")).Msg("update ticket comment failed: invalid comment id")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid comment id", "invalid comment id")
		return
	}

	userUUID, role, ok := getAuthenticatedUser(c)
	if !ok {
		log.Warn().Uint64("ticket_id", ticketID).Uint64("comment_id", commentID).Msg("update ticket comment failed: missing authenticated user context")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}

	var req requests.UpdateTicketCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Uint64("ticket_id", ticketID).Uint64("comment_id", commentID).Msg("update ticket comment failed: invalid request payload")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid request payload", "invalid request payload")
		return
	}

	comment, err := t.ticketService.UpdateComment(ticketID, commentID, userUUID, role, req.Body)
	if err != nil {
		if errors.Is(err, services.ErrTicketCommentForbidden) {
			log.Warn().Uint64("ticket_id", ticketID).Uint64("comment_id", commentID).Str("actor_role", string(role)).Msg("update ticket comment failed: forbidden")
			response.FailureWithAbort(c, http.StatusForbidden, "forbidden", "forbidden")
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Uint64("ticket_id", ticketID).Uint64("comment_id", commentID).Msg("update ticket comment failed: comment not found")
			response.FailureWithAbort(c, http.StatusNotFound, "comment not found", "comment not found")
			return
		}
		log.Error().Err(err).Uint64("ticket_id", ticketID).Uint64("comment_id", commentID).Msg("update ticket comment failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"comment_id":     comment.ID,
		"ticket_id":      comment.TicketID,
		"author_user_id": comment.AuthorUserID,
		"body":           comment.Body,
		"created_at":     comment.CreatedAt,
		"updated_at":     comment.UpdatedAt,
	}, "ticket comment updated")
}

func (t *TicketController) DeleteComment(c *gin.Context) {
	log := logger.L()
	// VULN-02: IDOR on tickets and comments — ticket id in path not access-controlled; comment delete is author/admin only.

	ticketID, ok := parseUintID(c.Param("id"))
	if !ok {
		log.Warn().Str("ticket_id", c.Param("id")).Msg("delete ticket comment failed: invalid ticket id")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid id", "invalid id")
		return
	}
	commentID, ok := parseUintID(c.Param("commentId"))
	if !ok {
		log.Warn().Str("comment_id", c.Param("commentId")).Msg("delete ticket comment failed: invalid comment id")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid comment id", "invalid comment id")
		return
	}

	userUUID, role, ok := getAuthenticatedUser(c)
	if !ok {
		log.Warn().Uint64("ticket_id", ticketID).Uint64("comment_id", commentID).Msg("delete ticket comment failed: missing authenticated user context")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}

	err := t.ticketService.DeleteComment(ticketID, commentID, userUUID, role)
	if err != nil {
		if errors.Is(err, services.ErrTicketCommentForbidden) {
			log.Warn().Uint64("ticket_id", ticketID).Uint64("comment_id", commentID).Str("actor_role", string(role)).Msg("delete ticket comment failed: forbidden")
			response.FailureWithAbort(c, http.StatusForbidden, "forbidden", "forbidden")
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Uint64("ticket_id", ticketID).Uint64("comment_id", commentID).Msg("delete ticket comment failed: comment not found")
			response.FailureWithAbort(c, http.StatusNotFound, "comment not found", "comment not found")
			return
		}
		log.Error().Err(err).Uint64("ticket_id", ticketID).Uint64("comment_id", commentID).Msg("delete ticket comment failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"comment_id": commentID, "ticket_id": ticketID}, "ticket comment deleted")
}

func parseUintID(raw string) (uint64, bool) {
	value, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, false
	}
	return value, true
}

func getAuthenticatedUser(c *gin.Context) (string, models.UserRole, bool) {
	userUUIDValue, ok := c.Get(middleware.CtxUserUUID)
	if !ok {
		return "", "", false
	}
	userUUID, ok := userUUIDValue.(string)
	if !ok || userUUID == "" {
		return "", "", false
	}

	roleValue, ok := c.Get(middleware.CtxUserRole)
	if !ok {
		return "", "", false
	}
	roleStr, ok := roleValue.(string)
	if !ok || roleStr == "" {
		return "", "", false
	}

	return userUUID, models.UserRole(roleStr), true
}

func collectUserIDsFromTickets(tickets []models.Ticket) []uint64 {
	seen := make(map[uint64]struct{})
	out := make([]uint64, 0)
	for _, tk := range tickets {
		if _, ok := seen[tk.ReporterUserID]; !ok {
			seen[tk.ReporterUserID] = struct{}{}
			out = append(out, tk.ReporterUserID)
		}
		if tk.AssignedUserID != nil {
			id := *tk.AssignedUserID
			if _, ok := seen[id]; !ok {
				seen[id] = struct{}{}
				out = append(out, id)
			}
		}
	}
	return out
}

func ticketDisplayName(u models.User) string {
	fn := strings.TrimSpace(u.FirstName)
	ln := strings.TrimSpace(u.LastName)
	if fn != "" || ln != "" {
		return strings.TrimSpace(fn + " " + ln)
	}
	return strings.TrimSpace(u.Email)
}

// formatTicketUserPublic returns a safe subset for API responses (no password hash).
func formatTicketUserPublic(u *models.User) interface{} {
	if u == nil {
		return nil
	}
	return gin.H{
		"user_id":      u.ID,
		"user_uuid":    u.UUID.String(),
		"email":        strings.TrimSpace(u.Email),
		"first_name":   strings.TrimSpace(u.FirstName),
		"last_name":    strings.TrimSpace(u.LastName),
		"display_name": ticketDisplayName(*u),
		"role":         u.Role,
	}
}

func (tc *TicketController) loadTicketUsers(tickets []models.Ticket) (map[uint64]models.User, error) {
	ids := collectUserIDsFromTickets(tickets)
	return tc.userRepo.GetMapByIDs(ids)
}

func formatTicket(ticket *models.Ticket, users map[uint64]models.User) gin.H {
	reporterPtr := ticket.Reporter
	if reporterPtr == nil {
		if u, ok := users[ticket.ReporterUserID]; ok {
			reporterPtr = &u
		}
	}

	var assigneePtr *models.User
	if ticket.Assignee != nil {
		assigneePtr = ticket.Assignee
	} else if ticket.AssignedUserID != nil {
		if u, ok := users[*ticket.AssignedUserID]; ok {
			assigneePtr = &u
		}
	}

	var assignedUserID any
	var assignedDisplay any
	var assignedEmail any
	if ticket.AssignedUserID != nil {
		assignedUserID = *ticket.AssignedUserID
		if assigneePtr != nil {
			assignedDisplay = ticketDisplayName(*assigneePtr)
			assignedEmail = strings.TrimSpace(assigneePtr.Email)
		}
	} else {
		assignedUserID = nil
		assignedDisplay = nil
		assignedEmail = nil
	}

	reporterDisplay := ""
	reporterEmail := ""
	if reporterPtr != nil {
		reporterDisplay = ticketDisplayName(*reporterPtr)
		reporterEmail = strings.TrimSpace(reporterPtr.Email)
	}

	return gin.H{
		"ticket_id":             ticket.ID,
		"reporter_user_id":      ticket.ReporterUserID,
		"reporter_display_name": reporterDisplay,
		"reporter_email":        reporterEmail,
		"assigned_user_id":      assignedUserID,
		"assigned_display_name": assignedDisplay,
		"assigned_email":        assignedEmail,
		"reporter":              formatTicketUserPublic(reporterPtr),
		"assignee":              formatTicketUserPublic(assigneePtr),
		"title":                 ticket.Title,
		"description":           ticket.Description,
		"category":              ticket.Category,
		"status":                ticket.Status,
		"created_at":            ticket.CreatedAt,
		"updated_at":            ticket.UpdatedAt,
	}
}
