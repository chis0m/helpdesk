package services

import (
	"errors"
	"strings"

	"github.com/google/uuid"

	"helpdesk/backend/internal/models"
	"helpdesk/backend/internal/repositories"
	"helpdesk/backend/internal/requests"
)

var ErrInvalidTicketStatusTransition = errors.New("invalid ticket status transition")

type TicketService struct {
	ticketRepo *repositories.TicketRepository
	userRepo   *repositories.UserRepository
}

func NewTicketService(ticketRepo *repositories.TicketRepository, userRepo *repositories.UserRepository) *TicketService {
	return &TicketService{
		ticketRepo: ticketRepo,
		userRepo:   userRepo,
	}
}

func (s *TicketService) CreateByUserUUID(userUUID string, req requests.CreateTicketRequest) (*models.Ticket, error) {
	parsedUUID, err := uuid.Parse(userUUID)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByUUID(parsedUUID)
	if err != nil {
		return nil, err
	}

	input := requests.CreateTicketInput{
		ReporterUserID: user.ID,
		Title:          strings.TrimSpace(req.Title),
		Description:    strings.TrimSpace(req.Description),
		Category:       strings.TrimSpace(req.Category),
	}
	return s.ticketRepo.Create(input)
}

func (s *TicketService) List(filter requests.ListTicketsFilter) ([]models.Ticket, int64, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	return s.ticketRepo.List(filter)
}

func (s *TicketService) GetByID(ticketID uint64) (*models.Ticket, error) {
	return s.ticketRepo.GetByID(ticketID)
}

func (s *TicketService) UpdateByID(ticketID uint64, req requests.UpdateTicketRequest) (*models.Ticket, error) {
	input := requests.UpdateTicketInput{
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
	}
	return s.ticketRepo.UpdateByID(ticketID, input)
}

func (s *TicketService) UpdateStatus(ticketID uint64, status models.TicketStatus) (*models.Ticket, error) {
	current, err := s.ticketRepo.GetByID(ticketID)
	if err != nil {
		return nil, err
	}
	if !isAllowedTransition(current.Status, status) {
		return nil, ErrInvalidTicketStatusTransition
	}
	return s.ticketRepo.UpdateStatus(ticketID, status)
}

func (s *TicketService) Assign(ticketID uint64, assignedUserID *uint64, unassign bool) (*models.Ticket, error) {
	if unassign {
		return s.ticketRepo.UpdateAssignment(ticketID, nil)
	}
	if assignedUserID == nil {
		return s.ticketRepo.GetByID(ticketID)
	}
	if _, err := s.userRepo.GetByID(*assignedUserID); err != nil {
		return nil, err
	}
	return s.ticketRepo.UpdateAssignment(ticketID, assignedUserID)
}

func (s *TicketService) DeleteByID(ticketID uint64) error {
	return s.ticketRepo.DeleteByID(ticketID)
}

func isAllowedTransition(from models.TicketStatus, to models.TicketStatus) bool {
	if from == to {
		return true
	}

	transitions := map[models.TicketStatus]map[models.TicketStatus]bool{
		models.TicketStatusOpen: {
			models.TicketStatusInProgress: true,
			models.TicketStatusClosed:     true,
		},
		models.TicketStatusInProgress: {
			models.TicketStatusOpen:     true,
			models.TicketStatusResolved: true,
			models.TicketStatusClosed:   true,
		},
		models.TicketStatusResolved: {
			models.TicketStatusInProgress: true,
			models.TicketStatusClosed:     true,
		},
		models.TicketStatusClosed: {
			models.TicketStatusOpen: true,
		},
	}

	return transitions[from][to]
}
