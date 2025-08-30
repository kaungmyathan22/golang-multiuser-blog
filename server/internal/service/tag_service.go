package service

import (
	"errors"
	"fmt"

	"github.com/kaungmyathan22/golang-multiuser-blog/internal/models"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/repository"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/utils"
)

type TagService interface {
	Create(req *models.TagCreateRequest) (*models.TagResponse, error)
	GetByID(id uint) (*models.TagResponse, error)
	GetBySlug(slug string) (*models.TagResponse, error)
	Update(tagID uint, req *models.TagUpdateRequest) (*models.TagResponse, error)
	Delete(tagID uint) error
	GetTags(page, perPage int) ([]models.TagResponse, models.PaginationMeta, error)
	GetAllTags() ([]models.TagResponse, error)
	GetPopularTags(limit int) ([]models.TagResponse, error)
}

type tagService struct {
	tagRepo repository.TagRepository
}

func NewTagService(tagRepo repository.TagRepository) TagService {
	return &tagService{
		tagRepo: tagRepo,
	}
}

func (s *tagService) Create(req *models.TagCreateRequest) (*models.TagResponse, error) {
	// Validate request
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, fmt.Errorf("validation failed: %v", validationErrors)
	}

	// Check if name is already taken
	if s.tagRepo.IsNameTaken(req.Name, 0) {
		return nil, errors.New("tag name is already taken")
	}

	// Generate slug from name
	slug := utils.GenerateSlug(req.Name)
	originalSlug := slug

	// Ensure slug is unique
	counter := 1
	for s.tagRepo.IsSlugTaken(slug, 0) {
		slug = fmt.Sprintf("%s-%d", originalSlug, counter)
		counter++
	}

	// Set default color if not provided
	color := req.Color
	if color == "" {
		color = "#3B82F6" // Default blue color
	}

	// Create tag
	tag := &models.Tag{
		Name:        utils.SanitizeText(req.Name),
		Slug:        slug,
		Description: utils.SanitizeText(req.Description),
		Color:       color,
	}

	if err := s.tagRepo.Create(tag); err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	response := tag.ToResponse()
	return &response, nil
}

func (s *tagService) GetByID(id uint) (*models.TagResponse, error) {
	tag, err := s.tagRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	response := tag.ToResponse()
	// Add post count
	response.PostsCount = len(tag.Posts)

	return &response, nil
}

func (s *tagService) GetBySlug(slug string) (*models.TagResponse, error) {
	tag, err := s.tagRepo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}

	response := tag.ToResponse()
	// Add post count
	response.PostsCount = len(tag.Posts)

	return &response, nil
}

func (s *tagService) Update(tagID uint, req *models.TagUpdateRequest) (*models.TagResponse, error) {
	// Validate request
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, fmt.Errorf("validation failed: %v", validationErrors)
	}

	// Get existing tag
	tag, err := s.tagRepo.GetByID(tagID)
	if err != nil {
		return nil, err
	}

	// Check if name is already taken (excluding current tag)
	if req.Name != "" && req.Name != tag.Name {
		if s.tagRepo.IsNameTaken(req.Name, tagID) {
			return nil, errors.New("tag name is already taken")
		}

		tag.Name = utils.SanitizeText(req.Name)

		// Regenerate slug if name changed
		newSlug := utils.GenerateSlug(req.Name)
		if newSlug != tag.Slug && !s.tagRepo.IsSlugTaken(newSlug, tagID) {
			tag.Slug = newSlug
		}
	}

	// Update other fields
	if req.Description != "" {
		tag.Description = utils.SanitizeText(req.Description)
	}

	if req.Color != "" {
		tag.Color = req.Color
	}

	if err := s.tagRepo.Update(tag); err != nil {
		return nil, fmt.Errorf("failed to update tag: %w", err)
	}

	response := tag.ToResponse()
	return &response, nil
}

func (s *tagService) Delete(tagID uint) error {
	// Check if tag exists
	_, err := s.tagRepo.GetByID(tagID)
	if err != nil {
		return err
	}

	return s.tagRepo.Delete(tagID)
}

func (s *tagService) GetTags(page, perPage int) ([]models.TagResponse, models.PaginationMeta, error) {
	offset := (page - 1) * perPage
	tags, total, err := s.tagRepo.List(offset, perPage)
	if err != nil {
		return nil, models.PaginationMeta{}, err
	}

	var responses []models.TagResponse
	for _, tag := range tags {
		responses = append(responses, tag.ToResponse())
	}

	pagination := utils.CalculatePagination(page, perPage, total)
	return responses, pagination, nil
}

func (s *tagService) GetAllTags() ([]models.TagResponse, error) {
	tags, err := s.tagRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []models.TagResponse
	for _, tag := range tags {
		responses = append(responses, tag.ToResponse())
	}

	return responses, nil
}

func (s *tagService) GetPopularTags(limit int) ([]models.TagResponse, error) {
	if limit <= 0 || limit > 50 {
		limit = 10 // Default limit
	}

	tags, err := s.tagRepo.GetPopular(limit)
	if err != nil {
		return nil, err
	}

	var responses []models.TagResponse
	for _, tag := range tags {
		response := tag.ToResponse()
		// Note: The posts_count is already calculated in the repository query
		responses = append(responses, response)
	}

	return responses, nil
}
