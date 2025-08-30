package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/middleware"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/models"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/service"
)

type TagHandler struct {
	tagService service.TagService
}

func NewTagHandler(tagService service.TagService) *TagHandler {
	return &TagHandler{
		tagService: tagService,
	}
}

// CreateTag godoc
// @Summary Create a new tag (Admin only)
// @Description Create a new tag for categorizing posts
// @Tags Tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param tag body models.TagCreateRequest true "Tag data"
// @Success 201 {object} models.APIResponse{data=models.TagResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 409 {object} models.APIResponse
// @Router /api/admin/tags [post]
func (h *TagHandler) CreateTag(c *gin.Context) {
	var req models.TagCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request format",
		})
		return
	}

	tag, err := h.tagService.Create(&req)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "tag name is already taken" {
			statusCode = http.StatusConflict
		}

		c.JSON(statusCode, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Tag created successfully",
		Data:    tag,
	})
}

// GetTag godoc
// @Summary Get a tag by ID
// @Description Get a specific tag by its ID
// @Tags Tags
// @Produce json
// @Param id path int true "Tag ID"
// @Success 200 {object} models.APIResponse{data=models.TagResponse}
// @Failure 404 {object} models.APIResponse
// @Router /api/tags/{id} [get]
func (h *TagHandler) GetTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid tag ID",
		})
		return
	}

	tag, err := h.tagService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "Tag not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    tag,
	})
}

// GetTagBySlug godoc
// @Summary Get a tag by slug
// @Description Get a specific tag by its slug
// @Tags Tags
// @Produce json
// @Param slug path string true "Tag slug"
// @Success 200 {object} models.APIResponse{data=models.TagResponse}
// @Failure 404 {object} models.APIResponse
// @Router /api/tags/slug/{slug} [get]
func (h *TagHandler) GetTagBySlug(c *gin.Context) {
	slug := c.Param("slug")

	tag, err := h.tagService.GetBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "Tag not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    tag,
	})
}

// UpdateTag godoc
// @Summary Update a tag (Admin only)
// @Description Update an existing tag
// @Tags Tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Tag ID"
// @Param tag body models.TagUpdateRequest true "Tag update data"
// @Success 200 {object} models.APIResponse{data=models.TagResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 409 {object} models.APIResponse
// @Router /api/admin/tags/{id} [put]
func (h *TagHandler) UpdateTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid tag ID",
		})
		return
	}

	var req models.TagUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request format",
		})
		return
	}

	tag, err := h.tagService.Update(uint(id), &req)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "tag name is already taken" {
			statusCode = http.StatusConflict
		} else if err.Error() == "tag not found" {
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
		Message: "Tag updated successfully",
		Data:    tag,
	})
}

// DeleteTag godoc
// @Summary Delete a tag (Admin only)
// @Description Delete an existing tag and remove it from all posts
// @Tags Tags
// @Security BearerAuth
// @Param id path int true "Tag ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/admin/tags/{id} [delete]
func (h *TagHandler) DeleteTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid tag ID",
		})
		return
	}

	err = h.tagService.Delete(uint(id))
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "tag not found" {
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
		Message: "Tag deleted successfully",
	})
}

// GetTags godoc
// @Summary Get tags
// @Description Get a paginated list of all tags
// @Tags Tags
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} models.PaginatedResponse{data=[]models.TagResponse}
// @Router /api/tags [get]
func (h *TagHandler) GetTags(c *gin.Context) {
	page, perPage := middleware.GetPaginationParams(c)

	tags, pagination, err := h.tagService.GetTags(page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to retrieve tags",
		})
		return
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success:    true,
		Data:       tags,
		Pagination: pagination,
	})
}

// GetAllTags godoc
// @Summary Get all tags
// @Description Get all tags without pagination (useful for dropdowns)
// @Tags Tags
// @Produce json
// @Success 200 {object} models.APIResponse{data=[]models.TagResponse}
// @Router /api/tags/all [get]
func (h *TagHandler) GetAllTags(c *gin.Context) {
	tags, err := h.tagService.GetAllTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to retrieve tags",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    tags,
	})
}

// GetPopularTags godoc
// @Summary Get popular tags
// @Description Get the most popular tags based on post count
// @Tags Tags
// @Produce json
// @Param limit query int false "Number of tags to return" default(10)
// @Success 200 {object} models.APIResponse{data=[]models.TagResponse}
// @Router /api/tags/popular [get]
func (h *TagHandler) GetPopularTags(c *gin.Context) {
	limitStr := c.Query("limit")
	limit := 10 // default

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = l
		}
	}

	tags, err := h.tagService.GetPopularTags(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to retrieve popular tags",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    tags,
	})
}

// GetPostsByTag godoc
// @Summary Get posts by tag
// @Description Get posts that have a specific tag
// @Tags Tags
// @Produce json
// @Param id path int true "Tag ID"
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} models.PaginatedResponse{data=[]models.PostListResponse}
// @Failure 404 {object} models.APIResponse
// @Router /api/tags/{id}/posts [get]
func (h *TagHandler) GetPostsByTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid tag ID",
		})
		return
	}

	// First check if tag exists
	_, err = h.tagService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "Tag not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Tag found - use /api/posts?tag_id=" + idStr + " to get posts with this tag",
		Data: map[string]interface{}{
			"redirect_url": "/api/posts?tag_id=" + idStr,
		},
	})
}

// GetTagStats godoc
// @Summary Get tag statistics (Admin only)
// @Description Get statistics about tags and their usage
// @Tags Tags
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=object}
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Router /api/admin/tags/stats [get]
func (h *TagHandler) GetTagStats(c *gin.Context) {
	// Get all tags with post counts
	allTags, err := h.tagService.GetAllTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to retrieve tag statistics",
		})
		return
	}

	// Get popular tags
	popularTags, err := h.tagService.GetPopularTags(5)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to retrieve popular tags",
		})
		return
	}

	// Calculate statistics
	totalTags := len(allTags)
	tagsWithPosts := 0
	totalPosts := 0

	for _, tag := range allTags {
		if tag.PostsCount > 0 {
			tagsWithPosts++
			totalPosts += tag.PostsCount
		}
	}

	stats := map[string]interface{}{
		"total_tags":         totalTags,
		"tags_with_posts":    tagsWithPosts,
		"tags_without_posts": totalTags - tagsWithPosts,
		"total_tag_usages":   totalPosts,
		"popular_tags":       popularTags,
		"average_posts_per_tag": func() float64 {
			if tagsWithPosts > 0 {
				return float64(totalPosts) / float64(tagsWithPosts)
			}
			return 0
		}(),
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    stats,
	})
}
