package controllers

import "github.com/gin-gonic/gin"

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (h *HealthController) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "ok",
		"message": "backend is running",
	})
}
