package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/url"
	"strings"
	"time"

	"gorm.io/gorm"

	"helpdesk/backend/internal/auth"
	"helpdesk/backend/internal/config"
	"helpdesk/backend/internal/mail"
	"helpdesk/backend/internal/models"
	"helpdesk/backend/internal/repositories"
	"helpdesk/backend/internal/requests"
)

var (
	ErrInviteNotFound    = errors.New("invite not found")
	ErrInviteExpired     = errors.New("invite expired")
	ErrInviteUsed        = errors.New("invite already used")
	ErrInviteEmailTaken  = errors.New("email already registered")
	ErrInviteForbidden   = errors.New("forbidden invite action")
	ErrInvitePendingExists = errors.New("pending invite already exists for this email")
)

type InviteService struct {
	cfg        config.Config
	inviteRepo *repositories.InviteRepository
	userRepo   *repositories.UserRepository
	notifier   mail.StaffInviteNotifier
}

func NewInviteService(
	cfg config.Config,
	inviteRepo *repositories.InviteRepository,
	userRepo *repositories.UserRepository,
	notifier mail.StaffInviteNotifier,
) *InviteService {
	return &InviteService{
		cfg:        cfg,
		inviteRepo: inviteRepo,
		userRepo:   userRepo,
		notifier:   notifier,
	}
}

func hashInviteToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func generateInviteRawToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func (s *InviteService) CreateStaffInvite(actorID uint64, actorRole models.UserRole, req requests.CreateStaffInviteRequest) (*models.Invite, error) {
	if actorRole != models.RoleAdmin && actorRole != models.RoleSuperAdmin {
		return nil, ErrInviteForbidden
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	now := time.Now().UTC()

	if _, err := s.userRepo.GetByEmail(email); err == nil {
		return nil, ErrInviteEmailTaken
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	pending, err := s.inviteRepo.HasPendingInviteForEmail(email, now)
	if err != nil {
		return nil, err
	}
	if pending {
		return nil, ErrInvitePendingExists
	}

	rawToken, err := generateInviteRawToken()
	if err != nil {
		return nil, err
	}
	tokenHash := hashInviteToken(rawToken)

	var middleName *string
	if req.MiddleName != nil {
		t := strings.TrimSpace(*req.MiddleName)
		if t != "" {
			middleName = &t
		}
	}

	inv := &models.Invite{
		Email:           email,
		TokenHash:       tokenHash,
		FirstName:       strings.TrimSpace(req.FirstName),
		LastName:        strings.TrimSpace(req.LastName),
		MiddleName:      middleName,
		InvitedByUserID: actorID,
		TargetRole:      models.RoleStaff,
		ExpiresAt:       now.Add(s.cfg.InviteTTL()),
	}

	if err := s.inviteRepo.Create(inv); err != nil {
		return nil, err
	}

	base := strings.TrimSuffix(strings.TrimSpace(s.cfg.FrontendURL), "/")
	inviteURL := base + "/accept-invite?token=" + url.QueryEscape(rawToken)
	if err := s.notifier.SendStaffInvite(email, inviteURL); err != nil {
		return nil, err
	}

	return inv, nil
}

type InviteVerifyResult struct {
	Valid     bool
	Email     string
	FirstName string
	LastName  string
}

func (s *InviteService) VerifyInvite(rawToken string) (*InviteVerifyResult, error) {
	rawToken = strings.TrimSpace(rawToken)
	if len(rawToken) != 64 {
		return &InviteVerifyResult{Valid: false}, nil
	}

	inv, err := s.inviteRepo.GetByTokenHash(hashInviteToken(rawToken))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &InviteVerifyResult{Valid: false}, nil
		}
		return nil, err
	}

	now := time.Now().UTC()
	if inv.UsedAt != nil {
		return &InviteVerifyResult{Valid: false}, nil
	}
	if !inv.ExpiresAt.UTC().After(now) {
		return &InviteVerifyResult{Valid: false}, nil
	}

	return &InviteVerifyResult{
		Valid:     true,
		Email:     inv.Email,
		FirstName: inv.FirstName,
		LastName:  inv.LastName,
	}, nil
}

func (s *InviteService) AcceptInvite(req requests.AcceptInviteRequest) (*models.User, error) {
	rawToken := strings.TrimSpace(req.Token)
	inv, err := s.inviteRepo.GetByTokenHash(hashInviteToken(rawToken))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInviteNotFound
		}
		return nil, err
	}

	now := time.Now().UTC()
	if inv.UsedAt != nil {
		return nil, ErrInviteUsed
	}
	if !inv.ExpiresAt.UTC().After(now) {
		return nil, ErrInviteExpired
	}

	if _, err := s.userRepo.GetByEmail(inv.Email); err == nil {
		return nil, ErrInviteEmailTaken
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	passwordHash, err := auth.HashPassword(strings.TrimSpace(req.Password))
	if err != nil {
		return nil, err
	}

	mustChange := false
	changedAt := now
	isActive := true
	createInput := requests.CreateUserInput{
		Email:              inv.Email,
		PasswordHash:       passwordHash,
		FirstName:          inv.FirstName,
		LastName:           inv.LastName,
		MiddleName:         inv.MiddleName,
		Role:               inv.TargetRole,
		IsActive:           &isActive,
		MustChangePassword: &mustChange,
		PasswordChangedAt:  &changedAt,
	}

	user, err := s.userRepo.Create(createInput)
	if err != nil {
		return nil, err
	}

	if err := s.inviteRepo.MarkUsed(inv.ID, now); err != nil {
		return nil, err
	}

	return user, nil
}
