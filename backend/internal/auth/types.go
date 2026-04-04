package auth

import "time"

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

const (
	AccessCookieName  = "access_token"
	RefreshCookieName = "refresh_token"
	CSRFHeaderName    = "X-CSRF-Token"
)

type TokenPair struct {
	AccessToken    string
	RefreshToken   string
	AccessExpires  time.Time
	RefreshExpires time.Time
	CSRFToken      string
	CSRFExpiresAt  time.Time
}

type CSRFToken struct {
	Token     string
	ExpiresAt time.Time
}
