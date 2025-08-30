package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/middleware"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/models"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/service"
)

type PostHandler struct {
	postService service.PostService
}

func NewPostHandler(postService service.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

// CreatePost godoc
// @Summary Create a new post
// @Description Create a new blog post
// @Tags Posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param post body models.PostCreateRequest true "Post data"
// @Success 201 {object} models.APIResponse{data=models.PostResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Router /api/posts [post]
func (h *PostHandler) CreatePost(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "User not authenticated",
		})
		return
	}

	var req models.PostCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request format",
		})
		return
	}

	post, err := h.postService.Create(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Post created successfully",
		Data:    post,
	})
}

// GetPost godoc
// @Summary Get a post by ID
// @Description Get a specific post by its ID
// @Tags Posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} models.APIResponse{data=models.PostResponse}
// @Failure 404 {object} models.APIResponse
// @Router /api/posts/{id} [get]
func (h *PostHandler) GetPost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid post ID",
		})
		return
	}

	post, err := h.postService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "Post not found",
		})
		return
	}

	// Increment view count for published posts
	if post.Status == models.PostStatusPublished {
		go h.postService.IncrementViewCount(uint(id))
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    post,
	})
}

// GetPostBySlug godoc
// @Summary Get a post by slug
// @Description Get a specific post by its slug
// @Tags Posts
// @Produce json
// @Param slug path string true "Post slug"
// @Success 200 {object} models.APIResponse{data=models.PostResponse}
// @Failure 404 {object} models.APIResponse
// @Router /api/posts/slug/{slug} [get]
func (h *PostHandler) GetPostBySlug(c *gin.Context) {
	slug := c.Param("slug")

	post, err := h.postService.GetBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "Post not found",
		})
		return
	}

	// Increment view count for published posts
	if post.Status == models.PostStatusPublished {
		go h.postService.IncrementViewCount(post.ID)
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    post,
	})
}

// UpdatePost godoc
// @Summary Update a post
// @Description Update an existing post
// @Tags Posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Post ID"
// @Param post body models.PostUpdateRequest true "Post update data"
// @Success 200 {object} models.APIResponse{data=models.PostResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/posts/{id} [put]
func (h *PostHandler) UpdatePost(c *gin.Context) {
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
			Error:   "Invalid post ID",
		})
		return
	}

	var req models.PostUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request format",
		})
		return
	}

	isAdmin := middleware.IsAdmin(c)
	post, err := h.postService.Update(uint(id), userID, &req, isAdmin)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "unauthorized: you can only update your own posts" {
			statusCode = http.StatusForbidden
		} else if err.Error() == "post not found" {
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
		Message: "Post updated successfully",
		Data:    post,
	})
}

// DeletePost godoc
// @Summary Delete a post
// @Description Delete an existing post
// @Tags Posts
// @Security BearerAuth
// @Param id path int true "Post ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/posts/{id} [delete]
func (h *PostHandler) DeletePost(c *gin.Context) {
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
			Error:   "Invalid post ID",
		})
		return
	}

	isAdmin := middleware.IsAdmin(c)
	err = h.postService.Delete(uint(id), userID, isAdmin)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "unauthorized: you can only delete your own posts" {
			statusCode = http.StatusForbidden
		} else if err.Error() == "post not found" {
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
		Message: "Post deleted successfully",
	})
}

// GetPosts godoc
// @Summary Get posts
// @Description Get a list of posts with pagination and filtering
// @Tags Posts
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param status query string false "Post status filter" Enums(draft, published, archived)
// @Param author_id query int false "Author ID filter"
// @Success 200 {object} models.PaginatedResponse{data=[]models.PostListResponse}
// @Router /api/posts [get]
func (h *PostHandler) GetPosts(c *gin.Context) {
	page, perPage := middleware.GetPaginationParams(c)

	var status models.PostStatus
	if statusStr := c.Query("status"); statusStr != "" {
		status = models.PostStatus(statusStr)
	}

	var authorID uint
	if authorIDStr := c.Query("author_id"); authorIDStr != "" {
		if id, err := strconv.ParseUint(authorIDStr, 10, 32); err == nil {
			authorID = uint(id)
		}
	}

	posts, pagination, err := h.postService.GetPosts(page, perPage, status, authorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to retrieve posts",
		})
		return
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success:    true,
		Data:       posts,
		Pagination: pagination,
	})
}

// GetPublishedPosts godoc
// @Summary Get published posts
// @Description Get a list of published posts
// @Tags Posts
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} models.PaginatedResponse{data=[]models.PostListResponse}
// @Router /api/posts/published [get]
func (h *PostHandler) GetPublishedPosts(c *gin.Context) {
	page, perPage := middleware.GetPaginationParams(c)

	posts, pagination, err := h.postService.GetPublishedPosts(page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to retrieve posts",
		})
		return
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success:    true,
		Data:       posts,
		Pagination: pagination,
	})
}

// SearchPosts godoc
// @Summary Search posts
// @Description Search for posts by title and content
// @Tags Posts
// @Produce json
// @Param q query string true "Search query"
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} models.PaginatedResponse{data=[]models.PostListResponse}
// @Router /api/posts/search [get]
func (h *PostHandler) SearchPosts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Search query is required",
		})
		return
	}

	page, perPage := middleware.GetPaginationParams(c)

	posts, pagination, err := h.postService.SearchPosts(query, page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to search posts",
		})
		return
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success:    true,
		Data:       posts,
		Pagination: pagination,
	})
}

// PublishPost godoc
// @Summary Publish a post
// @Description Publish a draft post
// @Tags Posts
// @Security BearerAuth
// @Param id path int true "Post ID"
// @Success 200 {object} models.APIResponse{data=models.PostResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/posts/{id}/publish [post]
func (h *PostHandler) PublishPost(c *gin.Context) {
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
			Error:   "Invalid post ID",
		})
		return
	}

	isAdmin := middleware.IsAdmin(c)
	post, err := h.postService.Publish(uint(id), userID, isAdmin)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "unauthorized: you can only publish your own posts" {
			statusCode = http.StatusForbidden
		} else if err.Error() == "post not found" {
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
		Message: "Post published successfully",
		Data:    post,
	})
}

// UnpublishPost godoc
// @Summary Unpublish a post
// @Description Unpublish a published post
// @Tags Posts
// @Security BearerAuth
// @Param id path int true "Post ID"
// @Success 200 {object} models.APIResponse{data=models.PostResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/posts/{id}/unpublish [post]
func (h *PostHandler) UnpublishPost(c *gin.Context) {
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
			Error:   "Invalid post ID",
		})
		return
	}

	isAdmin := middleware.IsAdmin(c)
	post, err := h.postService.Unpublish(uint(id), userID, isAdmin)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "unauthorized: you can only unpublish your own posts" {
			statusCode = http.StatusForbidden
		} else if err.Error() == "post not found" {
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
		Message: "Post unpublished successfully",
		Data:    post,
	})
}
