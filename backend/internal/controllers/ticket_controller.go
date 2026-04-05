package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"helpdesk/backend/internal/logger"
	"helpdesk/backend/internal/middleware"
	"helpdesk/backend/internal/models"
	"helpdesk/backend/internal/requests"
	"helpdesk/backend/internal/response"
	"helpdesk/backend/internal/services"
)

type TicketController struct {
	ticketService *services.TicketService
}

func NewTicketController(ticketService *services.TicketService) *TicketController {
	return &TicketController{ticketService: ticketService}
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

	response.Success(c, http.StatusCreated, formatTicket(ticket), "ticket created")
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

	tickets, total, err := t.ticketService.List(filter)
	if err != nil {
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

	result := make([]gin.H, 0, len(tickets))
	for _, ticket := range tickets {
		t := ticket
		result = append(result, formatTicket(&t))
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

func (t *TicketController) GetByID(c *gin.Context) {
	log := logger.L()

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

	response.Success(c, http.StatusOK, formatTicket(ticket), "ticket fetched")
}

func (t *TicketController) UpdateByID(c *gin.Context) {
	log := logger.L()

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

	response.Success(c, http.StatusOK, formatTicket(ticket), "ticket updated")
}

func (t *TicketController) UpdateStatus(c *gin.Context) {
	log := logger.L()

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

	response.Success(c, http.StatusOK, formatTicket(ticket), "ticket status updated")
}

func (t *TicketController) Assign(c *gin.Context) {
	log := logger.L()

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

	ticket, err := t.ticketService.Assign(ticketID, req.AssignedUserID, req.Unassign)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Uint64("ticket_id", ticketID).Msg("assign ticket failed: ticket or user not found")
			response.FailureWithAbort(c, http.StatusNotFound, "resource not found", "resource not found")
			return
		}
		log.Error().Err(err).Uint64("ticket_id", ticketID).Msg("assign ticket failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	response.Success(c, http.StatusOK, formatTicket(ticket), "ticket assignment updated")
}

func (t *TicketController) DeleteByID(c *gin.Context) {
	log := logger.L()

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

func parseUintID(raw string) (uint64, bool) {
	value, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, false
	}
	return value, true
}

func formatTicket(ticket *models.Ticket) gin.H {
	var assignedUserID any
	if ticket.AssignedUserID != nil {
		assignedUserID = *ticket.AssignedUserID
	}

	return gin.H{
		"ticket_id":        ticket.ID,
		"reporter_user_id": ticket.ReporterUserID,
		"assigned_user_id": assignedUserID,
		"title":            ticket.Title,
		"description":      ticket.Description,
		"category":         ticket.Category,
		"status":           ticket.Status,
		"created_at":       ticket.CreatedAt,
		"updated_at":       ticket.UpdatedAt,
	}
}
