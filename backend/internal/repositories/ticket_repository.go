package repositories

import (
	"errors"

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
		ReporterUserID: input.ReporterUserID,
		Title:          input.Title,
		Description:    input.Description,
		Category:       input.Category,
		Status:         models.TicketStatusOpen,
	}

	if err := r.db.Create(ticket).Error; err != nil {
		return nil, err
	}
	return ticket, nil
}

func (r *TicketRepository) GetByID(ticketID uint64) (*models.Ticket, error) {
	var ticket models.Ticket
	if err := r.db.First(&ticket, "id = ?", ticketID).Error; err != nil {
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

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	var tickets []models.Ticket
	if err := query.Order("created_at DESC").Offset(offset).Limit(filter.Limit).Find(&tickets).Error; err != nil {
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

	if err := r.db.First(&ticket, "id = ?", ticketID).Error; err != nil {
		return nil, errors.New("ticket updated but failed to reload")
	}
	return &ticket, nil
}

func (r *TicketRepository) UpdateStatus(ticketID uint64, status models.TicketStatus) (*models.Ticket, error) {
	var ticket models.Ticket
	if err := r.db.First(&ticket, "id = ?", ticketID).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&ticket).Update("status", status).Error; err != nil {
		return nil, err
	}

	if err := r.db.First(&ticket, "id = ?", ticketID).Error; err != nil {
		return nil, errors.New("ticket status updated but failed to reload")
	}
	return &ticket, nil
}

func (r *TicketRepository) UpdateAssignment(ticketID uint64, assignedUserID *uint64) (*models.Ticket, error) {
	var ticket models.Ticket
	if err := r.db.First(&ticket, "id = ?", ticketID).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&ticket).Update("assigned_user_id", assignedUserID).Error; err != nil {
		return nil, err
	}

	if err := r.db.First(&ticket, "id = ?", ticketID).Error; err != nil {
		return nil, errors.New("ticket assignment updated but failed to reload")
	}
	return &ticket, nil
}

func (r *TicketRepository) DeleteByID(ticketID uint64) error {
	return r.db.Delete(&models.Ticket{}, "id = ?", ticketID).Error
}
