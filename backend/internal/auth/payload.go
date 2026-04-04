package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Claims struct {
	Issuer    string        `json:"issuer"`
	Subject   string        `json:"subject"`
	Role      string        `json:"role"`
	Audience  string        `json:"audience"`
	Duration  time.Duration `json:"duration"`
	TokenType TokenType     `json:"token_type"`
	SessionID string        `json:"session_id"`
}

type Payload struct {
	Iss    string    `json:"iss"`
	Sub    string    `json:"sub"`
	Role   string    `json:"role"`
	Aud    string    `json:"aud"`
	Iat    time.Time `json:"iat"`
	Exp    time.Time `json:"exp"`
	Jti    string    `json:"jti"`
	Type   TokenType `json:"type"`
	SessID string    `json:"sess_id"`
}

func NewPayload(claims Claims) (*Payload, error) {
	if strings.TrimSpace(claims.Issuer) == "" {
		return nil, fmt.Errorf("issuer is required")
	}
	if strings.TrimSpace(claims.Subject) == "" {
		return nil, fmt.Errorf("subject is required")
	}
	if strings.TrimSpace(claims.Role) == "" {
		return nil, fmt.Errorf("role is required")
	}
	if strings.TrimSpace(claims.Audience) == "" {
		return nil, fmt.Errorf("audience is required")
	}
	if strings.TrimSpace(claims.SessionID) == "" {
		return nil, fmt.Errorf("session id is required")
	}
	if claims.Duration <= 0 {
		return nil, fmt.Errorf("token duration must be greater than zero")
	}

	now := time.Now().UTC()
	exp := now.Add(claims.Duration)
	jti, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate JTI: %w", err)
	}

	if claims.TokenType != TokenTypeAccess && claims.TokenType != TokenTypeRefresh {
		return nil, fmt.Errorf("invalid token type: %s", claims.TokenType)
	}

	payload := Payload{
		Iss:    claims.Issuer,
		Sub:    claims.Subject,
		Role:   claims.Role,
		Aud:    claims.Audience,
		Iat:    now,
		Exp:    exp,
		Jti:    jti.String(),
		Type:   claims.TokenType,
		SessID: claims.SessionID,
	}

	return &payload, nil
}

func (p *Payload) IsExpired() bool {
	return time.Now().UTC().After(p.Exp)
}

func (p *Payload) IsValidJTI(jti string) bool {
	return p.Jti == jti
}

func (p *Payload) IsValidSessionID(sessionID string) bool {
	return p.SessID == sessionID
}

func (p *Payload) IsValid(jti string, sessionID string) bool {
	return !p.IsExpired() && p.IsValidJTI(jti) && p.IsValidSessionID(sessionID)
}
