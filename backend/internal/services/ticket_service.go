package services

import (
	"errors"
	"strings"

	"github.com/google/uuid"

	"helpdesk/backend/internal/auth/policy"
	"helpdesk/backend/internal/models"
	"helpdesk/backend/internal/repositories"
	"helpdesk/backend/internal/requests"
)

var ErrInvalidTicketStatusTransition = errors.New("invalid ticket status transition")
var ErrTicketCommentForbidden = errors.New("forbidden comment action")
var ErrTicketListForbidden = errors.New("forbidden ticket list filters")
var ErrTicketAccessDenied = errors.New("ticket access denied")
var ErrTicketAssignForbidden = errors.New("forbidden ticket assignment")
var ErrTicketAssigneeInvalidRole = errors.New("assignee must be staff or admin")
var ErrTicketAlreadyAssignedToUser = errors.New("ticket already assigned to this user")
var ErrTicketAlreadyUnassigned = errors.New("ticket is not assigned")
var ErrTicketUnassignAssigneeMismatch = errors.New("assigned_user_id does not match ticket assignee")

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

func (s *TicketService) loadTicketForActor(ticketUUID uuid.UUID, actorUUID uuid.UUID, actorRole models.UserRole) (*models.Ticket, error) {
	ticket, err := s.ticketRepo.GetByUUID(ticketUUID)
	if err != nil {
		return nil, err
	}
	actor, err := s.userRepo.GetByUUID(actorUUID)
	if err != nil {
		return nil, err
	}
	if !policy.CanAccessTicket(actor, actorRole, ticket) {
		return nil, ErrTicketAccessDenied
	}
	return ticket, nil
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

func (s *TicketService) GetForActor(ticketUUID uuid.UUID, actorUUID uuid.UUID, actorRole models.UserRole) (*models.Ticket, error) {
	return s.loadTicketForActor(ticketUUID, actorUUID, actorRole)
}

func (s *TicketService) SearchTicketsSafe(keyword string) ([]models.Ticket, error) {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return []models.Ticket{}, nil
	}
	return s.ticketRepo.SearchByKeywordConcatSafe(keyword)
}

func (s *TicketService) SearchForActor(actorUUID uuid.UUID, actorRole models.UserRole, keyword string) ([]models.Ticket, error) {
	tickets, err := s.SearchTicketsSafe(keyword)
	if err != nil {
		return nil, err
	}
	actor, err := s.userRepo.GetByUUID(actorUUID)
	if err != nil {
		return nil, err
	}
	out := make([]models.Ticket, 0, len(tickets))
	for i := range tickets {
		tk := tickets[i]
		if policy.CanAccessTicket(actor, actorRole, &tk) {
			out = append(out, tk)
		}
	}
	return out, nil
}

func (s *TicketService) UpdateForActor(ticketUUID uuid.UUID, actorUUID uuid.UUID, actorRole models.UserRole, req requests.UpdateTicketRequest) (*models.Ticket, error) {
	ticket, err := s.loadTicketForActor(ticketUUID, actorUUID, actorRole)
	if err != nil {
		return nil, err
	}
	input := requests.UpdateTicketInput{}
	if req.Title != nil {
		t := strings.TrimSpace(*req.Title)
		input.Title = &t
	}
	if req.Description != nil {
		d := strings.TrimSpace(*req.Description)
		input.Description = &d
	}
	if req.Category != nil {
		cat := strings.TrimSpace(*req.Category)
		input.Category = &cat
	}
	return s.ticketRepo.UpdateByID(ticket.ID, input)
}

func (s *TicketService) UpdateStatusForActor(ticketUUID uuid.UUID, actorUUID uuid.UUID, actorRole models.UserRole, status models.TicketStatus) (*models.Ticket, error) {
	ticket, err := s.loadTicketForActor(ticketUUID, actorUUID, actorRole)
	if err != nil {
		return nil, err
	}
	if !isAllowedTransition(ticket.Status, status) {
		return nil, ErrInvalidTicketStatusTransition
	}
	return s.ticketRepo.UpdateStatus(ticket.ID, status)
}

func (s *TicketService) AssignForActor(ticketUUID uuid.UUID, actorUUID uuid.UUID, actorRole models.UserRole, assigneeUUID uuid.UUID, unassign bool) (*models.Ticket, error) {
	ticket, err := s.loadTicketForActor(ticketUUID, actorUUID, actorRole)
	if err != nil {
		return nil, err
	}
	if actorRole != models.RoleAdmin && actorRole != models.RoleSuperAdmin {
		return nil, ErrTicketAssignForbidden
	}
	assigneeUser, err := s.userRepo.GetByUUID(assigneeUUID)
	if err != nil {
		return nil, err
	}
	assigneeID := assigneeUser.ID
	if unassign {
		if ticket.AssignedUserID == nil {
			return nil, ErrTicketAlreadyUnassigned
		}
		if *ticket.AssignedUserID != assigneeID {
			return nil, ErrTicketUnassignAssigneeMismatch
		}
		return s.ticketRepo.UpdateAssignment(ticket.ID, nil)
	}
	if ticket.AssignedUserID != nil && *ticket.AssignedUserID == assigneeID {
		return nil, ErrTicketAlreadyAssignedToUser
	}
	if assigneeUser.Role != models.RoleStaff && assigneeUser.Role != models.RoleAdmin {
		return nil, ErrTicketAssigneeInvalidRole
	}
	return s.ticketRepo.UpdateAssignment(ticket.ID, &assigneeID)
}

func (s *TicketService) DeleteForActor(ticketUUID uuid.UUID, actorUUID uuid.UUID, actorRole models.UserRole) error {
	ticket, err := s.loadTicketForActor(ticketUUID, actorUUID, actorRole)
	if err != nil {
		return err
	}
	return s.ticketRepo.DeleteByID(ticket.ID)
}

func (s *TicketService) AddCommentForActor(ticketUUID uuid.UUID, actorUUID uuid.UUID, actorRole models.UserRole, body string) (*models.TicketComment, error) {
	ticket, err := s.loadTicketForActor(ticketUUID, actorUUID, actorRole)
	if err != nil {
		return nil, err
	}
	actor, err := s.userRepo.GetByUUID(actorUUID)
	if err != nil {
		return nil, err
	}
	input := requests.CreateTicketCommentInput{
		TicketID:     ticket.ID,
		AuthorUserID: actor.ID,
		Body:         strings.TrimSpace(body),
	}
	return s.commentRepo.Create(input)
}

func (s *TicketService) ListCommentsForActor(ticketUUID uuid.UUID, actorUUID uuid.UUID, actorRole models.UserRole) ([]models.TicketCommentWithAuthor, error) {
	ticket, err := s.loadTicketForActor(ticketUUID, actorUUID, actorRole)
	if err != nil {
		return nil, err
	}
	return s.commentRepo.ListByTicketID(ticket.ID)
}

func (s *TicketService) UpdateCommentForActor(ticketUUID uuid.UUID, actorUUID uuid.UUID, actorRole models.UserRole, commentID uint64, body string) (*models.TicketComment, error) {
	ticket, err := s.loadTicketForActor(ticketUUID, actorUUID, actorRole)
	if err != nil {
		return nil, err
	}
	return s.updateCommentBody(ticket.ID, commentID, actorUUID.String(), actorRole, body)
}

func (s *TicketService) DeleteCommentForActor(ticketUUID uuid.UUID, actorUUID uuid.UUID, actorRole models.UserRole, commentID uint64) error {
	ticket, err := s.loadTicketForActor(ticketUUID, actorUUID, actorRole)
	if err != nil {
		return err
	}
	return s.deleteComment(ticket.ID, commentID, actorUUID.String(), actorRole)
}

func (s *TicketService) updateCommentBody(ticketID uint64, commentID uint64, actorUserUUID string, role models.UserRole, body string) (*models.TicketComment, error) {
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

func (s *TicketService) deleteComment(ticketID uint64, commentID uint64, actorUserUUID string, role models.UserRole) error {
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
