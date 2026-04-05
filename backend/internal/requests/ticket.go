package requests

import "helpdesk/backend/internal/models"

type CreateTicketRequest struct {
	// VULN-03: Weak input validation / stored XSS risk — no HTML/script sanitization on ticket text.
	Title       string `json:"title" binding:"required,min=3,max=180"`
	Description string `json:"description" binding:"required,min=5,max=5000"`
	Category    string `json:"category" binding:"required,min=2,max=80"`
}

type CreateTicketInput struct {
	ReporterUserID uint64
	Title          string
	Description    string
	Category       string
}

type UpdateTicketRequest struct {
	// VULN-03: Weak input validation / stored XSS risk — no HTML/script sanitization on ticket text.
	Title       *string `json:"title" binding:"omitempty,min=3,max=180"`
	Description *string `json:"description" binding:"omitempty,min=5,max=5000"`
	Category    *string `json:"category" binding:"omitempty,min=2,max=80"`
}

type UpdateTicketInput struct {
	Title       *string
	Description *string
	Category    *string
}

type UpdateTicketStatusRequest struct {
	Status models.TicketStatus `json:"status" binding:"required,oneof=open in_progress resolved closed"`
}

type AssignTicketRequest struct {
	AssignedUserID *uint64 `json:"assigned_user_id" binding:"omitempty"`
	Unassign       bool    `json:"unassign" binding:"omitempty"`
}

type ListTicketsQuery struct {
	Page           int                  `form:"page"`
	Limit          int                  `form:"limit"`
	Status         *models.TicketStatus `form:"status" binding:"omitempty,oneof=open in_progress resolved closed"`
	Category       *string              `form:"category"`
	ReporterUserID *uint64              `form:"reporter_user_id"`
	AssignedUserID *uint64              `form:"assigned_user_id"`
}

type ListTicketsFilter struct {
	Page           int
	Limit          int
	Status         *models.TicketStatus
	Category       *string
	ReporterUserID *uint64
	AssignedUserID *uint64
	// ScopeToUserID limits rows to tickets the user reported or is assigned to (non-admin list).
	ScopeToUserID *uint64
}
