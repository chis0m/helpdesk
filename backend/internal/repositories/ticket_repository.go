package repositories

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"helpdesk/backend/internal/models"
	"helpdesk/backend/internal/requests"
)

type TicketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

func (r *TicketRepository) Create(input requests.CreateTicketInput) (*models.Ticket, error) {
	ticket := &models.Ticket{
		UUID:           uuid.New(),
		ReporterUserID: input.ReporterUserID,
		Title:          input.Title,
		Description:    input.Description,
		Category:       input.Category,
		Status:         models.TicketStatusOpen,
	}
	if err := r.db.Create(ticket).Error; err != nil {
		return nil, err
	}
	return r.GetByID(ticket.ID)
}

func (r *TicketRepository) GetByID(ticketID uint64) (*models.Ticket, error) {
	var ticket models.Ticket
	if err := r.db.Preload("Reporter").Preload("Assignee").First(&ticket, "id = ?", ticketID).Error; err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *TicketRepository) GetByUUID(ticketUUID uuid.UUID) (*models.Ticket, error) {
	var ticket models.Ticket
	if err := r.db.Preload("Reporter").Preload("Assignee").First(&ticket, "uuid = ?", ticketUUID).Error; err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *TicketRepository) List(filter requests.ListTicketsFilter) ([]models.Ticket, int64, error) {
	query := r.db.Model(&models.Ticket{})
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.Category != nil && *filter.Category != "" {
		query = query.Where("category = ?", *filter.Category)
	}
	if filter.ReporterUserID != nil {
		query = query.Where("reporter_user_id = ?", *filter.ReporterUserID)
	}
	if filter.AssignedUserID != nil {
		query = query.Where("assigned_user_id = ?", *filter.AssignedUserID)
	}
	if filter.ScopeToUserID != nil {
		uid := *filter.ScopeToUserID
		query = query.Where("(reporter_user_id = ? OR assigned_user_id = ?)", uid, uid)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	var tickets []models.Ticket
	if err := query.Preload("Reporter").Preload("Assignee").Order("created_at DESC").Offset(offset).Limit(filter.Limit).Find(&tickets).Error; err != nil {
		return nil, 0, err
	}
	return tickets, total, nil
}

func (r *TicketRepository) UpdateByID(ticketID uint64, input requests.UpdateTicketInput) (*models.Ticket, error) {
	var ticket models.Ticket
	if err := r.db.First(&ticket, "id = ?", ticketID).Error; err != nil {
		return nil, err
	}

	updates := map[string]any{}
	if input.Title != nil {
		updates["title"] = *input.Title
	}
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	if input.Category != nil {
		updates["category"] = *input.Category
	}

	if len(updates) == 0 {
		return &ticket, nil
	}

	if err := r.db.Model(&ticket).Updates(updates).Error; err != nil {
		return nil, err
	}

	return r.GetByID(ticketID)
}

func (r *TicketRepository) UpdateStatus(ticketID uint64, status models.TicketStatus) (*models.Ticket, error) {
	var ticket models.Ticket
	if err := r.db.First(&ticket, "id = ?", ticketID).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&ticket).Update("status", status).Error; err != nil {
		return nil, err
	}

	return r.GetByID(ticketID)
}

func (r *TicketRepository) UpdateAssignment(ticketID uint64, assignedUserID *uint64) (*models.Ticket, error) {
	var ticket models.Ticket
	if err := r.db.First(&ticket, "id = ?", ticketID).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&ticket).Update("assigned_user_id", assignedUserID).Error; err != nil {
		return nil, err
	}

	return r.GetByID(ticketID)
}

func (r *TicketRepository) DeleteByID(ticketID uint64) error {
	return r.db.Delete(&models.Ticket{}, "id = ?", ticketID).Error
}

// VULN-07: SQL injection (ticket keyword search) — q concatenated into Raw SQL without parameter binding.
func (r *TicketRepository) SearchByKeywordConcatUnsafe(q string) ([]models.Ticket, error) {
	//nolint:gosec // G201
	sqlStr := fmt.Sprintf(
		"SELECT * FROM tickets WHERE deleted_at IS NULL AND (title LIKE '%%%s%%' OR description LIKE '%%%s%%' OR category LIKE '%%%s%%') ORDER BY created_at DESC LIMIT 100",
		q, q, q,
	)
	var tickets []models.Ticket
	if err := r.db.Raw(sqlStr).Scan(&tickets).Error; err != nil {
		return nil, err
	}
	return tickets, nil
}
