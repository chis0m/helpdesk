package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"helpdesk/backend/internal/logger"
	"helpdesk/backend/internal/middleware"
	"helpdesk/backend/internal/models"
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

func (u *UserController) Create(c *gin.Context) {
	log := logger.L()

	var req requests.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("create user failed: invalid request payload")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid request payload", "invalid request payload")
		return
	}

	user, err := u.userService.CreateUserFromRequest(req)
	if err != nil {
		log.Error().Err(err).Msg("create user failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	response.Success(c, http.StatusCreated, gin.H{
		"user_id":   user.ID,
		"user_uuid": user.UUID.String(),
		"email":     user.Email,
		"role":      user.Role,
		"is_active": user.IsActive,
	}, "user created")
}

func (u *UserController) GetByID(c *gin.Context) {
	log := logger.L()

	userIDParam := c.Param("id")
	userID, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("user_id", userIDParam).Msg("get user failed: invalid id")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid id", "invalid id")
		return
	}

	user, err := u.userService.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Uint64("user_id", userID).Msg("get user failed: user not found")
			response.FailureWithAbort(c, http.StatusNotFound, "user not found", "user not found")
			return
		}
		log.Error().Err(err).Uint64("user_id", userID).Msg("get user failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"user_id":     user.ID,
		"user_uuid":   user.UUID.String(),
		"email":       user.Email,
		"first_name":  user.FirstName,
		"last_name":   user.LastName,
		"middle_name": user.MiddleName,
		"role":        user.Role,
		"is_active":   user.IsActive,
	}, "user fetched")
}

func (u *UserController) UpdateByID(c *gin.Context) {
	log := logger.L()

	userIDParam := c.Param("id")
	userID, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("user_id", userIDParam).Msg("update user failed: invalid id")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid id", "invalid id")
		return
	}

	var req requests.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("update user failed: invalid request payload")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid request payload", "invalid request payload")
		return
	}

	user, err := u.userService.UpdateByID(userID, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Uint64("user_id", userID).Msg("update user failed: user not found")
			response.FailureWithAbort(c, http.StatusNotFound, "user not found", "user not found")
			return
		}
		log.Error().Err(err).Uint64("user_id", userID).Msg("update user failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"user_id":     user.ID,
		"user_uuid":   user.UUID.String(),
		"email":       user.Email,
		"first_name":  user.FirstName,
		"last_name":   user.LastName,
		"middle_name": user.MiddleName,
		"role":        user.Role,
		"is_active":   user.IsActive,
	}, "user updated")
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
	actorRole := models.UserRole(roleStr)
	if !ok || (actorRole != models.RoleAdmin && actorRole != models.RoleSuperAdmin) {
		log.Warn().Str("role", roleStr).Msg("update role failed: admin or super_admin access required")
		response.FailureWithAbort(c, http.StatusForbidden, "admin or super_admin access required", "admin or super_admin access required")
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

	user, err := u.userService.UpdateRoleByIDAsActor(userID, req.Role, actorRole)
	if err != nil {
		if errors.Is(err, services.ErrUserRoleChangeForbidden) {
			log.Warn().Uint64("user_id", userID).Str("actor_role", string(actorRole)).Str("target_role", string(req.Role)).Msg("update role failed: forbidden role change")
			response.FailureWithAbort(c, http.StatusForbidden, "forbidden role change", "forbidden role change")
			return
		}
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

func (u *UserController) CreateStaff(c *gin.Context) {
	log := logger.L()

	roleValue, ok := c.Get(middleware.CtxUserRole)
	if !ok {
		log.Warn().Msg("create staff failed: missing user role in context")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}
	roleStr, ok := roleValue.(string)
	actorRole := models.UserRole(roleStr)
	if !ok || (actorRole != models.RoleAdmin && actorRole != models.RoleSuperAdmin) {
		log.Warn().Str("role", roleStr).Msg("create staff failed: admin or super_admin access required")
		response.FailureWithAbort(c, http.StatusForbidden, "admin or super_admin access required", "admin or super_admin access required")
		return
	}

	var req requests.CreateStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("create staff failed: invalid request payload")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid request payload", "invalid request payload")
		return
	}

	user, err := u.userService.CreateStaffFromRequest(req)
	if err != nil {
		log.Error().Err(err).Msg("create staff failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	response.Success(c, http.StatusCreated, gin.H{
		"user_id":   user.ID,
		"user_uuid": user.UUID.String(),
		"email":     user.Email,
		"role":      user.Role,
		"is_active": user.IsActive,
	}, "staff created")
}
