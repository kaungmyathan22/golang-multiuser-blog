package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/kaungmyathan22/golang-multiuser-blog/internal/models"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/repository"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/utils"
)

type PostService interface {
	Create(authorID uint, req *models.PostCreateRequest) (*models.PostResponse, error)
	GetByID(id uint) (*models.PostResponse, error)
	GetBySlug(slug string) (*models.PostResponse, error)
	Update(postID, authorID uint, req *models.PostUpdateRequest, isAdmin bool) (*models.PostResponse, error)
	Delete(postID, authorID uint, isAdmin bool) error
	GetPosts(page, perPage int, status models.PostStatus, authorID uint) ([]models.PostListResponse, models.PaginationMeta, error)
	GetPublishedPosts(page, perPage int) ([]models.PostListResponse, models.PaginationMeta, error)
	GetPostsByAuthor(authorID uint, page, perPage int) ([]models.PostListResponse, models.PaginationMeta, error)
	GetPostsByTag(tagID uint, page, perPage int) ([]models.PostListResponse, models.PaginationMeta, error)
	SearchPosts(query string, page, perPage int) ([]models.PostListResponse, models.PaginationMeta, error)
	IncrementViewCount(id uint) error
	Publish(postID, authorID uint, isAdmin bool) (*models.PostResponse, error)
	Unpublish(postID, authorID uint, isAdmin bool) (*models.PostResponse, error)
}

type postService struct {
	postRepo    repository.PostRepository
	tagRepo     repository.TagRepository
	commentRepo repository.CommentRepository
}

func NewPostService(postRepo repository.PostRepository, tagRepo repository.TagRepository, commentRepo repository.CommentRepository) PostService {
	return &postService{
		postRepo:    postRepo,
		tagRepo:     tagRepo,
		commentRepo: commentRepo,
	}
}

func (s *postService) Create(authorID uint, req *models.PostCreateRequest) (*models.PostResponse, error) {
	// Validate request
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, fmt.Errorf("validation failed: %v", validationErrors)
	}

	// Generate slug from title
	slug := utils.GenerateSlug(req.Title)
	originalSlug := slug

	// Ensure slug is unique
	counter := 1
	for s.postRepo.IsSlugTaken(slug, 0) {
		slug = fmt.Sprintf("%s-%d", originalSlug, counter)
		counter++
	}

	// Extract excerpt if not provided
	excerpt := req.Excerpt
	if excerpt == "" {
		excerpt = utils.ExtractExcerpt(req.Content, 200)
	}

	// Create post
	post := &models.Post{
		Title:       utils.SanitizeText(req.Title),
		Slug:        slug,
		Content:     req.Content,
		Excerpt:     utils.SanitizeText(excerpt),
		FeaturedImg: req.FeaturedImg,
		Status:      req.Status,
		AuthorID:    authorID,
	}

	// Set published date if status is published
	if req.Status == models.PostStatusPublished {
		now := time.Now()
		post.PublishedAt = &now
	}

	if err := s.postRepo.Create(post); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	// Add tags if provided
	if len(req.TagIDs) > 0 {
		if err := s.postRepo.UpdateTags(post.ID, req.TagIDs); err != nil {
			// Log error but don't fail the post creation
			fmt.Printf("Warning: Failed to add tags to post: %v\n", err)
		}
	}

	// Get the created post with relationships
	createdPost, err := s.postRepo.GetByID(post.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve created post: %w", err)
	}

	response := s.enrichPostResponse(createdPost)
	return &response, nil
}

func (s *postService) GetByID(id uint) (*models.PostResponse, error) {
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	response := s.enrichPostResponse(post)
	return &response, nil
}

func (s *postService) GetBySlug(slug string) (*models.PostResponse, error) {
	post, err := s.postRepo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}

	response := s.enrichPostResponse(post)
	return &response, nil
}

func (s *postService) Update(postID, authorID uint, req *models.PostUpdateRequest, isAdmin bool) (*models.PostResponse, error) {
	// Validate request
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, fmt.Errorf("validation failed: %v", validationErrors)
	}

	// Get existing post
	post, err := s.postRepo.GetByID(postID)
	if err != nil {
		return nil, err
	}

	// Check ownership (only author or admin can update)
	if !isAdmin && post.AuthorID != authorID {
		return nil, errors.New("unauthorized: you can only update your own posts")
	}

	// Update fields
	if req.Title != "" {
		post.Title = utils.SanitizeText(req.Title)

		// Regenerate slug if title changed
		newSlug := utils.GenerateSlug(req.Title)
		if newSlug != post.Slug && !s.postRepo.IsSlugTaken(newSlug, postID) {
			post.Slug = newSlug
		}
	}

	if req.Content != "" {
		post.Content = req.Content
	}

	if req.Excerpt != "" {
		post.Excerpt = utils.SanitizeText(req.Excerpt)
	} else if req.Content != "" {
		// Auto-generate excerpt from content
		post.Excerpt = utils.ExtractExcerpt(req.Content, 200)
	}

	if req.FeaturedImg != "" {
		post.FeaturedImg = req.FeaturedImg
	}

	// Handle status change
	if req.Status != "" && req.Status != post.Status {
		post.Status = req.Status

		// Set published date when publishing
		if req.Status == models.PostStatusPublished && post.PublishedAt == nil {
			now := time.Now()
			post.PublishedAt = &now
		}
	}

	if err := s.postRepo.Update(post); err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	// Update tags if provided
	if len(req.TagIDs) > 0 {
		if err := s.postRepo.UpdateTags(post.ID, req.TagIDs); err != nil {
			fmt.Printf("Warning: Failed to update tags for post: %v\n", err)
		}
	}

	// Get updated post with relationships
	updatedPost, err := s.postRepo.GetByID(post.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated post: %w", err)
	}

	response := s.enrichPostResponse(updatedPost)
	return &response, nil
}

func (s *postService) Delete(postID, authorID uint, isAdmin bool) error {
	// Get existing post
	post, err := s.postRepo.GetByID(postID)
	if err != nil {
		return err
	}

	// Check ownership (only author or admin can delete)
	if !isAdmin && post.AuthorID != authorID {
		return errors.New("unauthorized: you can only delete your own posts")
	}

	return s.postRepo.Delete(postID)
}

func (s *postService) GetPosts(page, perPage int, status models.PostStatus, authorID uint) ([]models.PostListResponse, models.PaginationMeta, error) {
	offset := (page - 1) * perPage
	posts, total, err := s.postRepo.List(offset, perPage, status, authorID)
	if err != nil {
		return nil, models.PaginationMeta{}, err
	}

	var responses []models.PostListResponse
	for _, post := range posts {
		response := s.enrichPostListResponse(&post)
		responses = append(responses, response)
	}

	pagination := utils.CalculatePagination(page, perPage, total)
	return responses, pagination, nil
}

func (s *postService) GetPublishedPosts(page, perPage int) ([]models.PostListResponse, models.PaginationMeta, error) {
	offset := (page - 1) * perPage
	posts, total, err := s.postRepo.GetPublished(offset, perPage)
	if err != nil {
		return nil, models.PaginationMeta{}, err
	}

	var responses []models.PostListResponse
	for _, post := range posts {
		response := s.enrichPostListResponse(&post)
		responses = append(responses, response)
	}

	pagination := utils.CalculatePagination(page, perPage, total)
	return responses, pagination, nil
}

func (s *postService) GetPostsByAuthor(authorID uint, page, perPage int) ([]models.PostListResponse, models.PaginationMeta, error) {
	offset := (page - 1) * perPage
	posts, total, err := s.postRepo.GetByAuthor(authorID, offset, perPage)
	if err != nil {
		return nil, models.PaginationMeta{}, err
	}

	var responses []models.PostListResponse
	for _, post := range posts {
		response := s.enrichPostListResponse(&post)
		responses = append(responses, response)
	}

	pagination := utils.CalculatePagination(page, perPage, total)
	return responses, pagination, nil
}

func (s *postService) GetPostsByTag(tagID uint, page, perPage int) ([]models.PostListResponse, models.PaginationMeta, error) {
	offset := (page - 1) * perPage
	posts, total, err := s.postRepo.GetByTag(tagID, offset, perPage)
	if err != nil {
		return nil, models.PaginationMeta{}, err
	}

	var responses []models.PostListResponse
	for _, post := range posts {
		response := s.enrichPostListResponse(&post)
		responses = append(responses, response)
	}

	pagination := utils.CalculatePagination(page, perPage, total)
	return responses, pagination, nil
}

func (s *postService) SearchPosts(query string, page, perPage int) ([]models.PostListResponse, models.PaginationMeta, error) {
	offset := (page - 1) * perPage
	posts, total, err := s.postRepo.Search(query, offset, perPage)
	if err != nil {
		return nil, models.PaginationMeta{}, err
	}

	var responses []models.PostListResponse
	for _, post := range posts {
		response := s.enrichPostListResponse(&post)
		responses = append(responses, response)
	}

	pagination := utils.CalculatePagination(page, perPage, total)
	return responses, pagination, nil
}

func (s *postService) IncrementViewCount(id uint) error {
	return s.postRepo.IncrementViewCount(id)
}

func (s *postService) Publish(postID, authorID uint, isAdmin bool) (*models.PostResponse, error) {
	post, err := s.postRepo.GetByID(postID)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if !isAdmin && post.AuthorID != authorID {
		return nil, errors.New("unauthorized: you can only publish your own posts")
	}

	post.Status = models.PostStatusPublished
	if post.PublishedAt == nil {
		now := time.Now()
		post.PublishedAt = &now
	}

	if err := s.postRepo.Update(post); err != nil {
		return nil, fmt.Errorf("failed to publish post: %w", err)
	}

	response := s.enrichPostResponse(post)
	return &response, nil
}

func (s *postService) Unpublish(postID, authorID uint, isAdmin bool) (*models.PostResponse, error) {
	post, err := s.postRepo.GetByID(postID)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if !isAdmin && post.AuthorID != authorID {
		return nil, errors.New("unauthorized: you can only unpublish your own posts")
	}

	post.Status = models.PostStatusDraft

	if err := s.postRepo.Update(post); err != nil {
		return nil, fmt.Errorf("failed to unpublish post: %w", err)
	}

	response := s.enrichPostResponse(post)
	return &response, nil
}

// Helper methods

func (s *postService) enrichPostResponse(post *models.Post) models.PostResponse {
	response := post.ToResponse()

	// Add tags
	var tagResponses []models.TagResponse
	for _, tag := range post.Tags {
		tagResponses = append(tagResponses, tag.ToResponse())
	}
	response.Tags = tagResponses

	// Add comment count
	commentCount, _ := s.commentRepo.CountByPost(post.ID)
	response.CommentsCount = int(commentCount)

	return response
}

func (s *postService) enrichPostListResponse(post *models.Post) models.PostListResponse {
	response := post.ToListResponse()

	// Add tags
	var tagResponses []models.TagResponse
	for _, tag := range post.Tags {
		tagResponses = append(tagResponses, tag.ToResponse())
	}
	response.Tags = tagResponses

	// Add comment count
	commentCount, _ := s.commentRepo.CountByPost(post.ID)
	response.CommentsCount = int(commentCount)

	return response
}
