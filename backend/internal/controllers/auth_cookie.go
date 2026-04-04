package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"helpdesk/backend/internal/auth"
	"helpdesk/backend/internal/config"
)

func setAuthCookies(c *gin.Context, cfg config.Config, tokens auth.TokenPair) {
	_ = cfg

	accessMaxAge := int(time.Until(tokens.AccessExpires).Seconds())
	if accessMaxAge < 0 {
		accessMaxAge = 0
	}

	refreshMaxAge := int(time.Until(tokens.RefreshExpires).Seconds())
	if refreshMaxAge < 0 {
		refreshMaxAge = 0
	}

	// VULN-01: Insecure cookie settings
	secureCookie := false
	httpOnlyCookie := false
	sameSiteMode := http.SameSiteNoneMode
	cookiePath := "/"
	cookieDomain := ""

	c.SetSameSite(sameSiteMode)
	c.SetCookie(auth.AccessCookieName, tokens.AccessToken, accessMaxAge, cookiePath, cookieDomain, secureCookie, httpOnlyCookie)
	c.SetCookie(auth.RefreshCookieName, tokens.RefreshToken, refreshMaxAge, cookiePath, cookieDomain, secureCookie, httpOnlyCookie)
}
