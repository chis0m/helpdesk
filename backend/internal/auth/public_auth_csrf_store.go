package auth

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrPublicAuthCSRFRequired = errors.New("public auth csrf token required")
	ErrPublicAuthCSRFInvalid  = errors.New("public auth csrf token invalid")
	ErrPublicAuthCSRFExpired  = errors.New("public auth csrf token expired")
)

// PublicAuthCSRFStore keeps short-lived auth CSRF tokens in memory.
// For production, I use Redis/shared storage for multi-instance deployments.
type PublicAuthCSRFStore struct {
	mu     sync.Mutex
	ttl    time.Duration
	tokens map[string]time.Time
}

func NewPublicAuthCSRFStore(ttl time.Duration) *PublicAuthCSRFStore {
	return &PublicAuthCSRFStore{
		ttl:    ttl,
		tokens: make(map[string]time.Time),
	}
}

func (s *PublicAuthCSRFStore) Issue() *CSRFToken {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.purgeExpiredLocked()

	token := uuid.NewString()
	expiresAt := time.Now().UTC().Add(s.ttl)
	s.tokens[token] = expiresAt

	return &CSRFToken{
		Token:     token,
		ExpiresAt: expiresAt,
	}
}

func (s *PublicAuthCSRFStore) ValidateAndConsume(token string) error {
	if token == "" {
		return ErrPublicAuthCSRFRequired
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	exp, ok := s.tokens[token]
	if !ok {
		return ErrPublicAuthCSRFInvalid
	}

	delete(s.tokens, token)

	if exp.UTC().Before(time.Now().UTC()) {
		return ErrPublicAuthCSRFExpired
	}

	return nil
}

func (s *PublicAuthCSRFStore) purgeExpiredLocked() {
	now := time.Now().UTC()
	for token, exp := range s.tokens {
		if exp.UTC().Before(now) {
			delete(s.tokens, token)
		}
	}
}
