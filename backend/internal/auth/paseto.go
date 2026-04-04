package auth

import (
	"errors"
	"fmt"

	"github.com/o1egl/paseto"
)

type MakerInterface interface {
	CreateToken(claims Claims) (string, *Payload, error)
	CreateAccessToken(claims Claims) (string, *Payload, error)
	CreateRefreshToken(claims Claims) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
	VerifyTokenType(token string, expectedType TokenType) (*Payload, error)
}

const symmetricKeySize = 32

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (*PasetoMaker, error) {
	if len(symmetricKey) != symmetricKeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", symmetricKeySize)
	}

	return &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}, nil
}

func (m *PasetoMaker) CreateToken(claims Claims) (string, *Payload, error) {
	payload, err := NewPayload(claims)
	if err != nil {
		return "", nil, err
	}

	token, err := m.paseto.Encrypt(m.symmetricKey, payload, nil)
	if err != nil {
		return "", nil, fmt.Errorf("failed to encrypt token: %w", err)
	}

	return token, payload, nil
}

func (m *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	if err := m.paseto.Decrypt(token, m.symmetricKey, payload, nil); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	if payload.IsExpired() {
		return nil, ErrExpiredToken
	}

	return payload, nil
}

func (m *PasetoMaker) CreateAccessToken(claims Claims) (string, *Payload, error) {
	claims.TokenType = TokenTypeAccess
	return m.CreateToken(claims)
}

func (m *PasetoMaker) CreateRefreshToken(claims Claims) (string, *Payload, error) {
	claims.TokenType = TokenTypeRefresh
	return m.CreateToken(claims)
}

func (m *PasetoMaker) VerifyTokenType(token string, expectedType TokenType) (*Payload, error) {
	payload, err := m.VerifyToken(token)
	if err != nil {
		return nil, err
	}

	if payload.Type != expectedType {
		return nil, fmt.Errorf("%w: expected token type %s", ErrInvalidToken, expectedType)
	}

	return payload, nil
}
