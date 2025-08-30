package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/middleware"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/models"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/service"
)

type AdminHandler struct {
	userService service.UserService
}

func NewAdminHandler(userService service.UserService) *AdminHandler {
	return &AdminHandler{
		userService: userService,
	}
}

// GetUsers godoc
// @Summary Get all users (Admin only)
// @Description Get a paginated list of all users
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} models.PaginatedResponse{data=[]models.UserResponse}
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Router /api/admin/users [get]
func (h *AdminHandler) GetUsers(c *gin.Context) {
	page, perPage := middleware.GetPaginationParams(c)

	users, pagination, err := h.userService.GetUsers(page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to retrieve users",
		})
		return
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success:    true,
		Data:       users,
		Pagination: pagination,
	})
}

// GetUser godoc
// @Summary Get user by ID (Admin only)
// @Description Get a specific user by their ID
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} models.APIResponse{data=models.UserResponse}
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/admin/users/{id} [get]
func (h *AdminHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    user,
	})
}

// DeactivateUser godoc
// @Summary Deactivate user (Admin only)
// @Description Deactivate a user account
// @Tags Admin
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/admin/users/{id}/deactivate [post]
func (h *AdminHandler) DeactivateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	// Prevent admin from deactivating themselves
	currentUserID, _ := middleware.GetUserID(c)
	if currentUserID == uint(id) {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "You cannot deactivate your own account",
		})
		return
	}

	err = h.userService.DeactivateUser(uint(id))
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "User deactivated successfully",
	})
}

// ActivateUser godoc
// @Summary Activate user (Admin only)
// @Description Activate a deactivated user account
// @Tags Admin
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/admin/users/{id}/activate [post]
func (h *AdminHandler) ActivateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	err = h.userService.ActivateUser(uint(id))
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "User activated successfully",
	})
}

// GetUserStats godoc
// @Summary Get user statistics (Admin only)
// @Description Get comprehensive statistics about users
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=object}
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Router /api/admin/users/stats [get]
func (h *AdminHandler) GetUserStats(c *gin.Context) {
	// Get all users to calculate statistics
	users, _, err := h.userService.GetUsers(1, 10000) // Get all users
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to retrieve user statistics",
		})
		return
	}

	// Calculate statistics
	var activeUsers, inactiveUsers, adminUsers, regularUsers int
	for _, user := range users {
		if user.IsActive {
			activeUsers++
		} else {
			inactiveUsers++
		}

		if user.IsAdmin {
			adminUsers++
		} else {
			regularUsers++
		}
	}

	stats := map[string]interface{}{
		"total_users":    len(users),
		"active_users":   activeUsers,
		"inactive_users": inactiveUsers,
		"admin_users":    adminUsers,
		"regular_users":  regularUsers,
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    stats,
	})
}

// GetDashboardStats godoc
// @Summary Get dashboard statistics (Admin only)
// @Description Get comprehensive statistics for the admin dashboard
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=object}
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Router /api/admin/dashboard/stats [get]
func (h *AdminHandler) GetDashboardStats(c *gin.Context) {
	// Note: This would typically involve multiple services
	// For now, we'll provide basic user statistics
	// In a full implementation, you'd also get post, comment, and tag statistics

	// Get user statistics
	users, _, err := h.userService.GetUsers(1, 10000)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to retrieve dashboard statistics",
		})
		return
	}

	// Calculate user statistics
	var activeUsers, totalUsers int
	totalUsers = len(users)
	for _, user := range users {
		if user.IsActive {
			activeUsers++
		}
	}

	stats := map[string]interface{}{
		"users": map[string]interface{}{
			"total":    totalUsers,
			"active":   activeUsers,
			"inactive": totalUsers - activeUsers,
		},
		// Note: In a real implementation, you would add:
		// "posts": { "total": totalPosts, "published": publishedPosts, "drafts": draftPosts },
		// "comments": { "total": totalComments, "pending": pendingComments, "approved": approvedComments },
		// "tags": { "total": totalTags, "used": usedTags },
		"last_updated": map[string]interface{}{
			"timestamp": "now", // In real implementation, use actual timestamps
		},
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    stats,
		Message: "Dashboard statistics retrieved successfully",
	})
}
