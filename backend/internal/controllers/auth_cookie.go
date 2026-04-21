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

	// VULN-01: Weak session cookie flags — HttpOnly/Secure false (SameSite=Strict for same-site dev; cookies still readable if HttpOnly false).
	secureCookie := false
	httpOnlyCookie := false
	cookiePath := "/"
	cookieDomain := strings.TrimSpace(cfg.CookieDomain)
	sameSite := http.SameSiteStrictMode

	c.SetSameSite(sameSite)
	c.SetCookie(auth.AccessCookieName, tokens.AccessToken, accessMaxAge, cookiePath, cookieDomain, secureCookie, httpOnlyCookie)
	c.SetCookie(auth.RefreshCookieName, tokens.RefreshToken, refreshMaxAge, cookiePath, cookieDomain, secureCookie, httpOnlyCookie)
}

func clearAuthCookies(c *gin.Context, cfg config.Config) {
	// VULN-01: Match setAuthCookies — same Path/Domain/SameSite so the browser clears the session cookies.
	secureCookie := false
	httpOnlyCookie := false
	cookiePath := "/"
	cookieDomain := strings.TrimSpace(cfg.CookieDomain)
	sameSite := http.SameSiteStrictMode

	c.SetSameSite(sameSite)
	c.SetCookie(auth.AccessCookieName, "", -1, cookiePath, cookieDomain, secureCookie, httpOnlyCookie)
	c.SetCookie(auth.RefreshCookieName, "", -1, cookiePath, cookieDomain, secureCookie, httpOnlyCookie)
}
