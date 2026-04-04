package auth

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var (
	ErrLoginCSRFRequired = errors.New("login csrf token required")
	ErrLoginCSRFInvalid  = errors.New("login csrf token invalid")
	ErrLoginCSRFExpired  = errors.New("login csrf token expired")
)

// LoginCSRFStore keeps short-lived login CSRF tokens in memory.
// Production we use redis for this.
type LoginCSRFStore struct {
	mu     sync.Mutex
	ttl    time.Duration
	tokens map[string]time.Time
}

func NewLoginCSRFStore(ttl time.Duration) *LoginCSRFStore {
	return &LoginCSRFStore{
		ttl:    ttl,
		tokens: make(map[string]time.Time),
	}
}

func (s *LoginCSRFStore) Issue() *CSRFToken {
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

func (s *LoginCSRFStore) ValidateAndConsume(token string) error {
	if token == "" {
		log.Warn().Msg("login csrf token is empty")
		return ErrLoginCSRFRequired
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	exp, ok := s.tokens[token]
	if !ok {
		log.Warn().Msg("login csrf token not found in memory")
		return ErrLoginCSRFInvalid
	}

	delete(s.tokens, token)

	if exp.UTC().Before(time.Now().UTC()) {
		log.Warn().Msg("login csrf token has expired")
		return ErrLoginCSRFExpired
	}

	return nil
}

func (s *LoginCSRFStore) purgeExpiredLocked() {
	now := time.Now().UTC()
	for token, exp := range s.tokens {
		if exp.UTC().Before(now) {
			delete(s.tokens, token)
		}
	}
}
