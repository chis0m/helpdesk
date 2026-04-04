package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"helpdesk/backend/internal/config"
)

const (
	accessCookieName  = "access_token"
	refreshCookieName = "refresh_token"
)

func setAuthCookies(c *gin.Context, cfg config.Config, accessToken, refreshToken string, accessExp, refreshExp time.Time) {
	_ = cfg

	accessMaxAge := int(time.Until(accessExp).Seconds())
	if accessMaxAge < 0 {
		accessMaxAge = 0
	}

	refreshMaxAge := int(time.Until(refreshExp).Seconds())
	if refreshMaxAge < 0 {
		refreshMaxAge = 0
	}

	// Baseline branch (intentionally weak cookie security):
	// - Secure=false
	// - HttpOnly=false
	// - SameSite=None
	secureCookie := false
	httpOnlyCookie := false
	cookiePath := "/"
	cookieDomain := ""

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(accessCookieName, accessToken, accessMaxAge, cookiePath, cookieDomain, secureCookie, httpOnlyCookie)
	c.SetCookie(refreshCookieName, refreshToken, refreshMaxAge, cookiePath, cookieDomain, secureCookie, httpOnlyCookie)
}
