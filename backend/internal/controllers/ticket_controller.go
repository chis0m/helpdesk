package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (t *TicketController) Search(c *gin.Context) {
	log := logger.L()

	q := strings.TrimSpace(c.Query("q"))
	if q == "" {
		response.FailureWithAbort(c, http.StatusBadRequest, "query parameter q is required", "query parameter q is required")
		return
	}

	userUUIDStr, role, ok := getAuthenticatedUser(c)
	if !ok {
		log.Warn().Msg("ticket search failed: missing authenticated user context")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}
	actorUUID, err := uuid.Parse(strings.TrimSpace(userUUIDStr))
	if err != nil {
		log.Warn().Err(err).Msg("ticket search failed: invalid actor uuid")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}

	tickets, err := t.ticketService.SearchForActor(actorUUID, role, q)
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

	ticketUUID, err := uuid.Parse(strings.TrimSpace(c.Param("id")))
	if err != nil {
		log.Warn().Str("ticket_uuid", c.Param("id")).Msg("get ticket failed: invalid uuid")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid ticket id", "invalid ticket id")
		return
	}

	userUUIDStr, role, ok := getAuthenticatedUser(c)
	if !ok {
		log.Warn().Msg("get ticket failed: missing authenticated user context")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}
	actorUUID, err := uuid.Parse(strings.TrimSpace(userUUIDStr))
	if err != nil {
		log.Warn().Err(err).Msg("get ticket failed: invalid actor uuid")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}

	ticket, err := t.ticketService.GetForActor(ticketUUID, actorUUID, role)
	if err != nil {
		if errors.Is(err, services.ErrTicketAccessDenied) {
			response.FailureWithAbort(c, http.StatusForbidden, "forbidden", "forbidden")
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Str("ticket_uuid", ticketUUID.String()).Msg("get ticket failed: ticket not found")
			response.FailureWithAbort(c, http.StatusNotFound, "ticket not found", "ticket not found")
			return
		}
		log.Error().Err(err).Str("ticket_uuid", ticketUUID.String()).Msg("get ticket failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	users, err := t.loadTicketUsers([]models.Ticket{*ticket})
	if err != nil {
		log.Error().Err(err).Str("ticket_uuid", ticketUUID.String()).Msg("get ticket failed: load user names")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}
	response.Success(c, http.StatusOK, formatTicket(ticket, users), "ticket fetched")
}

func (t *TicketController) UpdateByID(c *gin.Context) {
	log := logger.L()

	ticketUUID, err := uuid.Parse(strings.TrimSpace(c.Param("id")))
	if err != nil {
		log.Warn().Str("ticket_uuid", c.Param("id")).Msg("update ticket failed: invalid uuid")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid ticket id", "invalid ticket id")
		return
	}

	userUUIDStr, role, ok := getAuthenticatedUser(c)
	if !ok {
		log.Warn().Msg("update ticket failed: missing authenticated user context")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}
	actorUUID, err := uuid.Parse(strings.TrimSpace(userUUIDStr))
	if err != nil {
		log.Warn().Err(err).Msg("update ticket failed: invalid actor uuid")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}

	var req requests.UpdateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("update ticket failed: invalid request payload")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid request payload", "invalid request payload")
		return
	}

	ticket, err := t.ticketService.UpdateForActor(ticketUUID, actorUUID, role, req)
	if err != nil {
		if errors.Is(err, services.ErrTicketAccessDenied) {
			response.FailureWithAbort(c, http.StatusForbidden, "forbidden", "forbidden")
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Str("ticket_uuid", ticketUUID.String()).Msg("update ticket failed: ticket not found")
			response.FailureWithAbort(c, http.StatusNotFound, "ticket not found", "ticket not found")
			return
		}
		log.Error().Err(err).Str("ticket_uuid", ticketUUID.String()).Msg("update ticket failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	users, err := t.loadTicketUsers([]models.Ticket{*ticket})
	if err != nil {
		log.Error().Err(err).Str("ticket_uuid", ticketUUID.String()).Msg("update ticket failed: load user names")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}
	response.Success(c, http.StatusOK, formatTicket(ticket, users), "ticket updated")
}

func (t *TicketController) UpdateStatus(c *gin.Context) {
	log := logger.L()

	ticketUUID, err := uuid.Parse(strings.TrimSpace(c.Param("id")))
	if err != nil {
		log.Warn().Str("ticket_uuid", c.Param("id")).Msg("update ticket status failed: invalid uuid")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid ticket id", "invalid ticket id")
		return
	}

	userUUIDStr, role, ok := getAuthenticatedUser(c)
	if !ok {
		log.Warn().Msg("update ticket status failed: missing authenticated user context")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}
	actorUUID, err := uuid.Parse(strings.TrimSpace(userUUIDStr))
	if err != nil {
		log.Warn().Err(err).Msg("update ticket status failed: invalid actor uuid")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}

	var req requests.UpdateTicketStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("update ticket status failed: invalid request payload")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid request payload", "invalid request payload")
		return
	}

	ticket, err := t.ticketService.UpdateStatusForActor(ticketUUID, actorUUID, role, req.Status)
	if err != nil {
		if errors.Is(err, services.ErrInvalidTicketStatusTransition) {
			log.Warn().Str("ticket_uuid", ticketUUID.String()).Str("requested_status", string(req.Status)).Msg("update ticket status failed: invalid status transition")
			response.FailureWithAbort(c, http.StatusBadRequest, "invalid status transition", "invalid status transition")
			return
		}
		if errors.Is(err, services.ErrTicketAccessDenied) {
			response.FailureWithAbort(c, http.StatusForbidden, "forbidden", "forbidden")
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Str("ticket_uuid", ticketUUID.String()).Msg("update ticket status failed: ticket not found")
			response.FailureWithAbort(c, http.StatusNotFound, "ticket not found", "ticket not found")
			return
		}
		log.Error().Err(err).Str("ticket_uuid", ticketUUID.String()).Msg("update ticket status failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	users, err := t.loadTicketUsers([]models.Ticket{*ticket})
	if err != nil {
		log.Error().Err(err).Str("ticket_uuid", ticketUUID.String()).Msg("update ticket status failed: load user names")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}
	response.Success(c, http.StatusOK, formatTicket(ticket, users), "ticket status updated")
}

func (t *TicketController) Assign(c *gin.Context) {
	log := logger.L()

	ticketUUID, err := uuid.Parse(strings.TrimSpace(c.Param("id")))
	if err != nil {
		log.Warn().Str("ticket_uuid", c.Param("id")).Msg("assign ticket failed: invalid uuid")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid ticket id", "invalid ticket id")
		return
	}

	userUUIDStr, role, ok := getAuthenticatedUser(c)
	if !ok {
		log.Warn().Msg("assign ticket failed: missing authenticated user context")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}
	actorUUID, err := uuid.Parse(strings.TrimSpace(userUUIDStr))
	if err != nil {
		log.Warn().Err(err).Msg("assign ticket failed: invalid actor uuid")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}

	var req requests.AssignTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("assign ticket failed: invalid request payload")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid request payload", "invalid request payload")
		return
	}

	ticket, err := t.ticketService.AssignForActor(ticketUUID, actorUUID, role, req.AssignedUserID, req.Unassign)
	if err != nil {
		if errors.Is(err, services.ErrTicketAccessDenied) {
			response.FailureWithAbort(c, http.StatusForbidden, "forbidden", "forbidden")
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Str("ticket_uuid", ticketUUID.String()).Msg("assign ticket failed: ticket or user not found")
			response.FailureWithAbort(c, http.StatusNotFound, "resource not found", "resource not found")
			return
		}
		log.Error().Err(err).Str("ticket_uuid", ticketUUID.String()).Msg("assign ticket failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	users, err := t.loadTicketUsers([]models.Ticket{*ticket})
	if err != nil {
		log.Error().Err(err).Str("ticket_uuid", ticketUUID.String()).Msg("assign ticket failed: load user names")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}
	response.Success(c, http.StatusOK, formatTicket(ticket, users), "ticket assignment updated")
}

func (t *TicketController) DeleteByID(c *gin.Context) {
	log := logger.L()

	ticketUUID, err := uuid.Parse(strings.TrimSpace(c.Param("id")))
	if err != nil {
		log.Warn().Str("ticket_uuid", c.Param("id")).Msg("delete ticket failed: invalid uuid")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid ticket id", "invalid ticket id")
		return
	}

	userUUIDStr, role, ok := getAuthenticatedUser(c)
	if !ok {
		log.Warn().Msg("delete ticket failed: missing authenticated user context")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}
	actorUUID, err := uuid.Parse(strings.TrimSpace(userUUIDStr))
	if err != nil {
		log.Warn().Err(err).Msg("delete ticket failed: invalid actor uuid")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}

	if err := t.ticketService.DeleteForActor(ticketUUID, actorUUID, role); err != nil {
		if errors.Is(err, services.ErrTicketAccessDenied) {
			response.FailureWithAbort(c, http.StatusForbidden, "forbidden", "forbidden")
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Str("ticket_uuid", ticketUUID.String()).Msg("delete ticket failed: ticket not found")
			response.FailureWithAbort(c, http.StatusNotFound, "ticket not found", "ticket not found")
			return
		}
		log.Error().Err(err).Str("ticket_uuid", ticketUUID.String()).Msg("delete ticket failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"ticket_uuid": ticketUUID.String()}, "ticket deleted")
}

func (t *TicketController) AddComment(c *gin.Context) {
	log := logger.L()

	ticketUUID, err := uuid.Parse(strings.TrimSpace(c.Param("id")))
	if err != nil {
		log.Warn().Str("ticket_uuid", c.Param("id")).Msg("add ticket comment failed: invalid uuid")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid ticket id", "invalid ticket id")
		return
	}

	userUUIDStr, role, ok := getAuthenticatedUser(c)
	if !ok {
		log.Warn().Msg("add ticket comment failed: missing authenticated user context")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}
	actorUUID, err := uuid.Parse(strings.TrimSpace(userUUIDStr))
	if err != nil {
		log.Warn().Err(err).Msg("add ticket comment failed: invalid actor uuid")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}

	var req requests.CreateTicketCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("add ticket comment failed: invalid request payload")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid request payload", "invalid request payload")
		return
	}

	comment, err := t.ticketService.AddCommentForActor(ticketUUID, actorUUID, role, req.Body)
	if err != nil {
		if errors.Is(err, services.ErrTicketAccessDenied) {
			response.FailureWithAbort(c, http.StatusForbidden, "forbidden", "forbidden")
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Str("ticket_uuid", ticketUUID.String()).Msg("add ticket comment failed: ticket or actor not found")
			response.FailureWithAbort(c, http.StatusNotFound, "resource not found", "resource not found")
			return
		}
		log.Error().Err(err).Str("ticket_uuid", ticketUUID.String()).Msg("add ticket comment failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	response.Success(c, http.StatusCreated, gin.H{
		"comment_id":     comment.ID,
		"comment_uuid":   comment.UUID.String(),
		"ticket_id":      comment.TicketID,
		"author_user_id": comment.AuthorUserID,
		"body":           comment.Body,
		"created_at":     comment.CreatedAt,
		"updated_at":     comment.UpdatedAt,
	}, "ticket comment created")
}

func (t *TicketController) ListComments(c *gin.Context) {
	log := logger.L()

	ticketUUID, err := uuid.Parse(strings.TrimSpace(c.Param("id")))
	if err != nil {
		log.Warn().Str("ticket_uuid", c.Param("id")).Msg("list ticket comments failed: invalid uuid")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid ticket id", "invalid ticket id")
		return
	}

	userUUIDStr, role, ok := getAuthenticatedUser(c)
	if !ok {
		log.Warn().Msg("list ticket comments failed: missing authenticated user context")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}
	actorUUID, err := uuid.Parse(strings.TrimSpace(userUUIDStr))
	if err != nil {
		log.Warn().Err(err).Msg("list ticket comments failed: invalid actor uuid")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}

	comments, err := t.ticketService.ListCommentsForActor(ticketUUID, actorUUID, role)
	if err != nil {
		if errors.Is(err, services.ErrTicketAccessDenied) {
			response.FailureWithAbort(c, http.StatusForbidden, "forbidden", "forbidden")
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Str("ticket_uuid", ticketUUID.String()).Msg("list ticket comments failed: ticket not found")
			response.FailureWithAbort(c, http.StatusNotFound, "ticket not found", "ticket not found")
			return
		}
		log.Error().Err(err).Str("ticket_uuid", ticketUUID.String()).Msg("list ticket comments failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	if comments == nil {
		comments = []models.TicketCommentWithAuthor{}
	}

	response.Success(c, http.StatusOK, gin.H{"items": comments}, "ticket comments fetched")
}

func (t *TicketController) UpdateComment(c *gin.Context) {
	log := logger.L()

	ticketUUID, err := uuid.Parse(strings.TrimSpace(c.Param("id")))
	if err != nil {
		log.Warn().Str("ticket_uuid", c.Param("id")).Msg("update ticket comment failed: invalid uuid")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid ticket id", "invalid ticket id")
		return
	}
	commentID, ok := parseUintID(c.Param("commentId"))
	if !ok {
		log.Warn().Str("comment_id", c.Param("commentId")).Msg("update ticket comment failed: invalid comment id")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid comment id", "invalid comment id")
		return
	}

	userUUIDStr, role, ok := getAuthenticatedUser(c)
	if !ok {
		log.Warn().Msg("update ticket comment failed: missing authenticated user context")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}
	actorUUID, err := uuid.Parse(strings.TrimSpace(userUUIDStr))
	if err != nil {
		log.Warn().Err(err).Msg("update ticket comment failed: invalid actor uuid")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}

	var req requests.UpdateTicketCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("update ticket comment failed: invalid request payload")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid request payload", "invalid request payload")
		return
	}

	comment, err := t.ticketService.UpdateCommentForActor(ticketUUID, actorUUID, role, commentID, req.Body)
	if err != nil {
		if errors.Is(err, services.ErrTicketAccessDenied) {
			response.FailureWithAbort(c, http.StatusForbidden, "forbidden", "forbidden")
			return
		}
		if errors.Is(err, services.ErrTicketCommentForbidden) {
			log.Warn().Str("ticket_uuid", ticketUUID.String()).Uint64("comment_id", commentID).Str("actor_role", string(role)).Msg("update ticket comment failed: forbidden")
			response.FailureWithAbort(c, http.StatusForbidden, "forbidden", "forbidden")
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Str("ticket_uuid", ticketUUID.String()).Uint64("comment_id", commentID).Msg("update ticket comment failed: comment not found")
			response.FailureWithAbort(c, http.StatusNotFound, "comment not found", "comment not found")
			return
		}
		log.Error().Err(err).Str("ticket_uuid", ticketUUID.String()).Uint64("comment_id", commentID).Msg("update ticket comment failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"comment_id":     comment.ID,
		"comment_uuid":   comment.UUID.String(),
		"ticket_id":      comment.TicketID,
		"author_user_id": comment.AuthorUserID,
		"body":           comment.Body,
		"created_at":     comment.CreatedAt,
		"updated_at":     comment.UpdatedAt,
	}, "ticket comment updated")
}

func (t *TicketController) DeleteComment(c *gin.Context) {
	log := logger.L()

	ticketUUID, err := uuid.Parse(strings.TrimSpace(c.Param("id")))
	if err != nil {
		log.Warn().Str("ticket_uuid", c.Param("id")).Msg("delete ticket comment failed: invalid uuid")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid ticket id", "invalid ticket id")
		return
	}
	commentID, ok := parseUintID(c.Param("commentId"))
	if !ok {
		log.Warn().Str("comment_id", c.Param("commentId")).Msg("delete ticket comment failed: invalid comment id")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid comment id", "invalid comment id")
		return
	}

	userUUIDStr, role, ok := getAuthenticatedUser(c)
	if !ok {
		log.Warn().Msg("delete ticket comment failed: missing authenticated user context")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}
	actorUUID, err := uuid.Parse(strings.TrimSpace(userUUIDStr))
	if err != nil {
		log.Warn().Err(err).Msg("delete ticket comment failed: invalid actor uuid")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}

	err = t.ticketService.DeleteCommentForActor(ticketUUID, actorUUID, role, commentID)
	if err != nil {
		if errors.Is(err, services.ErrTicketAccessDenied) {
			response.FailureWithAbort(c, http.StatusForbidden, "forbidden", "forbidden")
			return
		}
		if errors.Is(err, services.ErrTicketCommentForbidden) {
			log.Warn().Str("ticket_uuid", ticketUUID.String()).Uint64("comment_id", commentID).Str("actor_role", string(role)).Msg("delete ticket comment failed: forbidden")
			response.FailureWithAbort(c, http.StatusForbidden, "forbidden", "forbidden")
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Str("ticket_uuid", ticketUUID.String()).Uint64("comment_id", commentID).Msg("delete ticket comment failed: comment not found")
			response.FailureWithAbort(c, http.StatusNotFound, "comment not found", "comment not found")
			return
		}
		log.Error().Err(err).Str("ticket_uuid", ticketUUID.String()).Uint64("comment_id", commentID).Msg("delete ticket comment failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"comment_id": commentID, "ticket_uuid": ticketUUID.String()}, "ticket comment deleted")
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

// SEC-02: Public ticket JSON uses opaque ticket_uuid only; internal numeric id is not exposed.
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
		"ticket_uuid":           ticket.UUID.String(),
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
