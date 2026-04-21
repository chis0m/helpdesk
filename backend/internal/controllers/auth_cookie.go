package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"helpdesk/backend/internal/auth"
	"helpdesk/backend/internal/config"
)

func setAuthCookies(c *gin.Context, cfg config.Config, tokens auth.TokenPair) {
	accessMaxAge := int(time.Until(tokens.AccessExpires).Seconds())
	if accessMaxAge < 0 {
		accessMaxAge = 0
	}

	refreshMaxAge := int(time.Until(tokens.RefreshExpires).Seconds())
	if refreshMaxAge < 0 {
		refreshMaxAge = 0
	}

	// SECURE-01: HttpOnly always; Secure when production or COOKIE_SECURE=true (plain HTTP localhost keeps Secure=false).
	secureCookie := cfg.CookieSecure()
	httpOnlyCookie := true
	cookiePath := "/"
	cookieDomain := strings.TrimSpace(cfg.CookieDomain)
	sameSite := http.SameSiteStrictMode

	c.SetSameSite(sameSite)
	c.SetCookie(auth.AccessCookieName, tokens.AccessToken, accessMaxAge, cookiePath, cookieDomain, secureCookie, httpOnlyCookie)
	c.SetCookie(auth.RefreshCookieName, tokens.RefreshToken, refreshMaxAge, cookiePath, cookieDomain, secureCookie, httpOnlyCookie)
}

func clearAuthCookies(c *gin.Context, cfg config.Config) {
	// SECURE-01: Match setAuthCookies Path/Domain/SameSite/Secure/HttpOnly so logout clears cookies reliably.
	secureCookie := cfg.CookieSecure()
	httpOnlyCookie := true
	cookiePath := "/"
	cookieDomain := strings.TrimSpace(cfg.CookieDomain)
	sameSite := http.SameSiteStrictMode

	c.SetSameSite(sameSite)
	c.SetCookie(auth.AccessCookieName, "", -1, cookiePath, cookieDomain, secureCookie, httpOnlyCookie)
	c.SetCookie(auth.RefreshCookieName, "", -1, cookiePath, cookieDomain, secureCookie, httpOnlyCookie)
}
