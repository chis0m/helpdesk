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
var ErrTicketCommentForbidden = errors.New("forbidden comment action")
var ErrTicketListForbidden = errors.New("forbidden ticket list filters")

type TicketService struct {
	ticketRepo  *repositories.TicketRepository
	commentRepo *repositories.TicketCommentRepository
	userRepo    *repositories.UserRepository
}

func NewTicketService(
	ticketRepo *repositories.TicketRepository,
	commentRepo *repositories.TicketCommentRepository,
	userRepo *repositories.UserRepository,
) *TicketService {
	return &TicketService{
		ticketRepo:  ticketRepo,
		commentRepo: commentRepo,
		userRepo:    userRepo,
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

// ListForActor applies authorization: only admin and super_admin may list without scope;
// other roles only see tickets they reported or are assigned to.
func (s *TicketService) ListForActor(actorUserUUID string, actorRole models.UserRole, filter requests.ListTicketsFilter) ([]models.Ticket, int64, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}

	parsedUUID, err := uuid.Parse(strings.TrimSpace(actorUserUUID))
	if err != nil {
		return nil, 0, err
	}
	actor, err := s.userRepo.GetByUUID(parsedUUID)
	if err != nil {
		return nil, 0, err
	}

	if actorRole == models.RoleAdmin || actorRole == models.RoleSuperAdmin {
		filter.ScopeToUserID = nil
		return s.ticketRepo.List(filter)
	}

	if filter.ReporterUserID != nil && *filter.ReporterUserID != actor.ID {
		return nil, 0, ErrTicketListForbidden
	}
	if filter.AssignedUserID != nil && *filter.AssignedUserID != actor.ID {
		return nil, 0, ErrTicketListForbidden
	}

	scope := actor.ID
	filter.ScopeToUserID = &scope
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

func (s *TicketService) AddComment(ticketID uint64, actorUserUUID string, body string) (*models.TicketComment, error) {
	if _, err := s.ticketRepo.GetByID(ticketID); err != nil {
		return nil, err
	}

	actor, err := s.getActorByUUID(actorUserUUID)
	if err != nil {
		return nil, err
	}

	input := requests.CreateTicketCommentInput{
		TicketID:     ticketID,
		AuthorUserID: actor.ID,
		Body:         strings.TrimSpace(body),
	}
	return s.commentRepo.Create(input)
}

func (s *TicketService) ListComments(ticketID uint64) ([]models.TicketCommentWithAuthor, error) {
	if _, err := s.ticketRepo.GetByID(ticketID); err != nil {
		return nil, err
	}
	return s.commentRepo.ListByTicketID(ticketID)
}

func (s *TicketService) UpdateComment(ticketID uint64, commentID uint64, actorUserUUID string, role models.UserRole, body string) (*models.TicketComment, error) {
	comment, err := s.commentRepo.GetByIDAndTicketID(commentID, ticketID)
	if err != nil {
		return nil, err
	}

	actor, err := s.getActorByUUID(actorUserUUID)
	if err != nil {
		return nil, err
	}

	if role != models.RoleAdmin && role != models.RoleSuperAdmin && comment.AuthorUserID != actor.ID {
		return nil, ErrTicketCommentForbidden
	}

	return s.commentRepo.UpdateBody(commentID, strings.TrimSpace(body))
}

func (s *TicketService) DeleteComment(ticketID uint64, commentID uint64, actorUserUUID string, role models.UserRole) error {
	comment, err := s.commentRepo.GetByIDAndTicketID(commentID, ticketID)
	if err != nil {
		return err
	}

	actor, err := s.getActorByUUID(actorUserUUID)
	if err != nil {
		return err
	}

	if role != models.RoleAdmin && role != models.RoleSuperAdmin && comment.AuthorUserID != actor.ID {
		return ErrTicketCommentForbidden
	}

	return s.commentRepo.DeleteByID(commentID)
}

func (s *TicketService) getActorByUUID(actorUserUUID string) (*models.User, error) {
	parsedUUID, err := uuid.Parse(actorUserUUID)
	if err != nil {
		return nil, err
	}
	return s.userRepo.GetByUUID(parsedUUID)
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
