package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"helpdesk/backend/internal/config"
	"helpdesk/backend/internal/response"
)

type AppDetailController struct {
	cfg config.Config
}

func NewAppDetailController(cfg config.Config) *AppDetailController {
	return &AppDetailController{cfg: cfg}
}

type appDetailResponse struct {
	AppName     string `json:"app_name"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

// Get returns public app metadata (unauthenticated). Used by the SPA landing and auth entry pages.
func (a *AppDetailController) Get(c *gin.Context) {
	name := strings.TrimSpace(a.cfg.AppName)
	if name == "" {
		name = "SecWeb HelpDesk"
	}
	env := strings.TrimSpace(a.cfg.GoEnv)
	if env == "" {
		env = "development"
	}
	ver := strings.TrimSpace(a.cfg.AppVersion)
	if ver == "" {
		ver = "v1.0.0"
	}
	response.Success(c, http.StatusOK, appDetailResponse{
		AppName:     name,
		Environment: env,
		Version:     ver,
	}, "app detail fetched successfully")
}
