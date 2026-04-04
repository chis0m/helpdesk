package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"helpdesk/backend/internal/logger"
	"helpdesk/backend/internal/models"
	"helpdesk/backend/internal/middleware"
	"helpdesk/backend/internal/requests"
	"helpdesk/backend/internal/response"
	"helpdesk/backend/internal/services"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{userService: userService}
}

func (u *UserController) UpdateRoleByUserID(c *gin.Context) {
	log := logger.L()

	roleValue, ok := c.Get(middleware.CtxUserRole)
	if !ok {
		log.Warn().Msg("update role failed: missing user role in context")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}
	roleStr, ok := roleValue.(string)
	if !ok || models.UserRole(roleStr) != models.RoleAdmin {
		log.Warn().Str("role", roleStr).Msg("update role failed: admin access required")
		response.FailureWithAbort(c, http.StatusForbidden, "admin access required", "admin access required")
		return
	}

	userIDParam := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("user_id", userIDParam).Msg("update role failed: invalid user_id")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid user_id", "invalid user_id")
		return
	}

	var req requests.UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("update role failed: invalid request payload")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid request payload", "invalid request payload")
		return
	}

	user, err := u.userService.UpdateRoleByID(userID, req.Role)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Uint64("user_id", userID).Msg("update role failed: user not found")
			response.FailureWithAbort(c, http.StatusNotFound, "user not found", "user not found")
			return
		}
		log.Error().Err(err).Uint64("user_id", userID).Msg("update role failed: server error")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"user_id":   user.ID,
		"user_uuid": user.UUID.String(),
		"role":      user.Role,
	}, "user role updated")
}
