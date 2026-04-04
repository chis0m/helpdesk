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
)

type TokenPair struct {
	AccessToken    string
	RefreshToken   string
	AccessExpires  time.Time
	RefreshExpires time.Time
}
