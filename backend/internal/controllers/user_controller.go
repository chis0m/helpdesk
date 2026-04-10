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

// VULN-02: IDOR on user profiles — GET /users/:id without checking path id matches authenticated user.
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

// VULN-02: IDOR on user profiles — PATCH /users/:id without checking path id matches authenticated user.
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
	if !ok || actorRole != models.RoleSuperAdmin {
		log.Warn().Str("role", roleStr).Msg("update role failed: super_admin access required")
		response.FailureWithAbort(c, http.StatusForbidden, "super_admin access required", "super_admin access required")
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

func (u *UserController) ListAdmin(c *gin.Context) {
	log := logger.L()

	roleValue, ok := c.Get(middleware.CtxUserRole)
	if !ok {
		log.Warn().Msg("list users failed: missing user role in context")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}
	roleStr, ok := roleValue.(string)
	if !ok || roleStr == "" {
		log.Warn().Msg("list users failed: invalid role in context")
		response.FailureWithAbort(c, http.StatusUnauthorized, "authentication required", "authentication required")
		return
	}
	actorRole := models.UserRole(roleStr)
	if actorRole != models.RoleAdmin && actorRole != models.RoleSuperAdmin {
		log.Warn().Str("role", roleStr).Msg("list users failed: admin or super_admin access required")
		response.FailureWithAbort(c, http.StatusForbidden, "admin or super_admin access required", "admin or super_admin access required")
		return
	}

	var query requests.ListUsersQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		log.Warn().Err(err).Msg("list users failed: invalid query params")
		response.FailureWithAbort(c, http.StatusBadRequest, "invalid query parameters", "invalid query parameters")
		return
	}

	page := query.Page
	if page < 1 {
		page = 1
	}
	limit := query.Limit
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	var roleFilter *models.UserRole
	if query.Role != "" {
		r := models.UserRole(query.Role)
		roleFilter = &r
	}

	users, total, err := u.userService.ListAll(page, limit, roleFilter)
	if err != nil {
		log.Error().Err(err).Msg("list users failed")
		response.FailureWithAbort(c, http.StatusInternalServerError, "internal server error", "internal server error")
		return
	}

	items := make([]gin.H, 0, len(users))
	for i := range users {
		uu := users[i]
		items = append(items, gin.H{
			"user_id":     uu.ID,
			"user_uuid":   uu.UUID.String(),
			"email":       uu.Email,
			"first_name":  uu.FirstName,
			"last_name":   uu.LastName,
			"middle_name": uu.MiddleName,
			"role":        uu.Role,
			"is_active":   uu.IsActive,
			"created_at":  uu.CreatedAt,
			"updated_at":  uu.UpdatedAt,
		})
	}

	response.Success(c, http.StatusOK, gin.H{
		"items": items,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	}, "users fetched")
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

	user, err := u.userService.CreateStaffFromRequest(actorRole, req)
	if err != nil {
		if errors.Is(err, services.ErrCreateStaffAdminForbidden) {
			log.Warn().Str("actor_role", string(actorRole)).Msg("create staff failed: only admin or super_admin may set role admin")
			response.FailureWithAbort(c, http.StatusForbidden, "only admin or super_admin may create staff with role admin", "only admin or super_admin may create staff with role admin")
			return
		}
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
