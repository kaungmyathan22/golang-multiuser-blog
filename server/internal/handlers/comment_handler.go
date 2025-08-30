package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/middleware"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/models"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/service"
)

type CommentHandler struct {
	commentService service.CommentService
}

func NewCommentHandler(commentService service.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

// CreateComment godoc
// @Summary Create a new comment
// @Description Create a new comment or reply to an existing comment
// @Tags Comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param comment body models.CommentCreateRequest true "Comment data"
// @Success 201 {object} models.APIResponse{data=models.CommentResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Router /api/comments [post]
func (h *CommentHandler) CreateComment(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "User not authenticated",
		})
		return
	}

	var req models.CommentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request format",
		})
		return
	}

	comment, err := h.commentService.Create(userID, &req)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "post not found" || err.Error() == "parent comment not found" {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Comment created successfully (pending approval)",
		Data:    comment,
	})
}

// GetComment godoc
// @Summary Get a comment by ID
// @Description Get a specific comment by its ID
// @Tags Comments
// @Produce json
// @Param id path int true "Comment ID"
// @Success 200 {object} models.APIResponse{data=models.CommentResponse}
// @Failure 404 {object} models.APIResponse
// @Router /api/comments/{id} [get]
func (h *CommentHandler) GetComment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid comment ID",
		})
		return
	}

	comment, err := h.commentService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "Comment not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    comment,
	})
}

// UpdateComment godoc
// @Summary Update a comment
// @Description Update an existing comment
// @Tags Comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Comment ID"
// @Param comment body models.CommentUpdateRequest true "Comment update data"
// @Success 200 {object} models.APIResponse{data=models.CommentResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/comments/{id} [put]
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "User not authenticated",
		})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid comment ID",
		})
		return
	}

	var req models.CommentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request format",
		})
		return
	}

	isAdmin := middleware.IsAdmin(c)
	comment, err := h.commentService.Update(uint(id), userID, &req, isAdmin)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "unauthorized: you can only update your own comments" {
			statusCode = http.StatusForbidden
		} else if err.Error() == "comment not found" {
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
		Message: "Comment updated successfully",
		Data:    comment,
	})
}

// DeleteComment godoc
// @Summary Delete a comment
// @Description Delete an existing comment and its replies
// @Tags Comments
// @Security BearerAuth
// @Param id path int true "Comment ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/comments/{id} [delete]
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "User not authenticated",
		})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid comment ID",
		})
		return
	}

	isAdmin := middleware.IsAdmin(c)
	err = h.commentService.Delete(uint(id), userID, isAdmin)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "unauthorized: you can only delete your own comments" {
			statusCode = http.StatusForbidden
		} else if err.Error() == "comment not found" {
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
		Message: "Comment deleted successfully",
	})
}

// GetCommentsByPost godoc
// @Summary Get comments for a post
// @Description Get paginated comments for a specific post
// @Tags Comments
// @Produce json
// @Param post_id path int true "Post ID"
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} models.PaginatedResponse{data=[]models.CommentResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/posts/{post_id}/comments [get]
func (h *CommentHandler) GetCommentsByPost(c *gin.Context) {
	postIDStr := c.Param("post_id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid post ID",
		})
		return
	}

	page, perPage := middleware.GetPaginationParams(c)

	comments, pagination, err := h.commentService.GetByPost(uint(postID), page, perPage)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "post not found" {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success:    true,
		Data:       comments,
		Pagination: pagination,
	})
}

// GetCommentsByAuthor godoc
// @Summary Get comments by author
// @Description Get paginated comments by a specific author
// @Tags Comments
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} models.PaginatedResponse{data=[]models.CommentResponse}
// @Failure 401 {object} models.APIResponse
// @Router /api/comments/my-comments [get]
func (h *CommentHandler) GetCommentsByAuthor(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "User not authenticated",
		})
		return
	}

	page, perPage := middleware.GetPaginationParams(c)

	comments, pagination, err := h.commentService.GetByAuthor(userID, page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to retrieve comments",
		})
		return
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success:    true,
		Data:       comments,
		Pagination: pagination,
	})
}

// GetPendingComments godoc
// @Summary Get pending comments (Admin only)
// @Description Get paginated list of comments pending approval
// @Tags Comments
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} models.PaginatedResponse{data=[]models.CommentResponse}
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Router /api/admin/comments/pending [get]
func (h *CommentHandler) GetPendingComments(c *gin.Context) {
	page, perPage := middleware.GetPaginationParams(c)

	comments, pagination, err := h.commentService.GetPending(page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to retrieve pending comments",
		})
		return
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success:    true,
		Data:       comments,
		Pagination: pagination,
	})
}

// ApproveComment godoc
// @Summary Approve a comment (Admin only)
// @Description Approve a pending comment
// @Tags Comments
// @Security BearerAuth
// @Param id path int true "Comment ID"
// @Success 200 {object} models.APIResponse{data=models.CommentResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/admin/comments/{id}/approve [post]
func (h *CommentHandler) ApproveComment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid comment ID",
		})
		return
	}

	comment, err := h.commentService.ApproveComment(uint(id))
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "comment not found" {
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
		Message: "Comment approved successfully",
		Data:    comment,
	})
}

// RejectComment godoc
// @Summary Reject a comment (Admin only)
// @Description Reject a pending comment
// @Tags Comments
// @Security BearerAuth
// @Param id path int true "Comment ID"
// @Success 200 {object} models.APIResponse{data=models.CommentResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/admin/comments/{id}/reject [post]
func (h *CommentHandler) RejectComment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid comment ID",
		})
		return
	}

	comment, err := h.commentService.RejectComment(uint(id))
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "comment not found" {
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
		Message: "Comment rejected successfully",
		Data:    comment,
	})
}

// GetPendingCount godoc
// @Summary Get pending comments count (Admin only)
// @Description Get the total number of comments pending approval
// @Tags Comments
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=object{count=int}}
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Router /api/admin/comments/pending/count [get]
func (h *CommentHandler) GetPendingCount(c *gin.Context) {
	count, err := h.commentService.GetPendingCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to get pending comments count",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"count": count,
		},
	})
}
