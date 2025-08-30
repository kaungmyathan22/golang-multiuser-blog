package service

import (
	"errors"
	"fmt"

	"github.com/kaungmyathan22/golang-multiuser-blog/internal/models"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/repository"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/utils"
)

type CommentService interface {
	Create(authorID uint, req *models.CommentCreateRequest) (*models.CommentResponse, error)
	GetByID(id uint) (*models.CommentResponse, error)
	Update(commentID, authorID uint, req *models.CommentUpdateRequest, isAdmin bool) (*models.CommentResponse, error)
	Delete(commentID, authorID uint, isAdmin bool) error
	GetByPost(postID uint, page, perPage int) ([]models.CommentResponse, models.PaginationMeta, error)
	GetByAuthor(authorID uint, page, perPage int) ([]models.CommentResponse, models.PaginationMeta, error)
	GetPending(page, perPage int) ([]models.CommentResponse, models.PaginationMeta, error)
	ApproveComment(commentID uint) (*models.CommentResponse, error)
	RejectComment(commentID uint) (*models.CommentResponse, error)
	GetPendingCount() (int64, error)
}

type commentService struct {
	commentRepo repository.CommentRepository
	postRepo    repository.PostRepository
}

func NewCommentService(commentRepo repository.CommentRepository, postRepo repository.PostRepository) CommentService {
	return &commentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
	}
}

func (s *commentService) Create(authorID uint, req *models.CommentCreateRequest) (*models.CommentResponse, error) {
	// Validate request
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, fmt.Errorf("validation failed: %v", validationErrors)
	}

	// Verify that the post exists
	_, err := s.postRepo.GetByID(req.PostID)
	if err != nil {
		return nil, errors.New("post not found")
	}

	// Verify parent comment exists if this is a reply
	if req.ParentID != nil {
		_, err := s.commentRepo.GetByID(*req.ParentID)
		if err != nil {
			return nil, errors.New("parent comment not found")
		}
	}

	// Create comment
	comment := &models.Comment{
		Content:  utils.SanitizeText(req.Content),
		AuthorID: authorID,
		PostID:   req.PostID,
		ParentID: req.ParentID,
		Status:   models.CommentStatusPending, // Comments need approval by default
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	// Get the created comment with relationships
	createdComment, err := s.commentRepo.GetByID(comment.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve created comment: %w", err)
	}

	response := createdComment.ToResponse()
	return &response, nil
}

func (s *commentService) GetByID(id uint) (*models.CommentResponse, error) {
	comment, err := s.commentRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	response := comment.ToResponse()
	return &response, nil
}

func (s *commentService) Update(commentID, authorID uint, req *models.CommentUpdateRequest, isAdmin bool) (*models.CommentResponse, error) {
	// Validate request
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, fmt.Errorf("validation failed: %v", validationErrors)
	}

	// Get existing comment
	comment, err := s.commentRepo.GetByID(commentID)
	if err != nil {
		return nil, err
	}

	// Check ownership (only author or admin can update)
	if !isAdmin && comment.AuthorID != authorID {
		return nil, errors.New("unauthorized: you can only update your own comments")
	}

	// Update fields
	if req.Content != "" {
		comment.Content = utils.SanitizeText(req.Content)
		// Reset status to pending if content is changed (except by admin)
		if !isAdmin {
			comment.Status = models.CommentStatusPending
		}
	}

	// Only admin can change status directly
	if isAdmin && req.Status != "" {
		comment.Status = req.Status
	}

	if err := s.commentRepo.Update(comment); err != nil {
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}

	response := comment.ToResponse()
	return &response, nil
}

func (s *commentService) Delete(commentID, authorID uint, isAdmin bool) error {
	// Get existing comment
	comment, err := s.commentRepo.GetByID(commentID)
	if err != nil {
		return err
	}

	// Check ownership (only author or admin can delete)
	if !isAdmin && comment.AuthorID != authorID {
		return errors.New("unauthorized: you can only delete your own comments")
	}

	return s.commentRepo.Delete(commentID)
}

func (s *commentService) GetByPost(postID uint, page, perPage int) ([]models.CommentResponse, models.PaginationMeta, error) {
	// Verify that the post exists
	_, err := s.postRepo.GetByID(postID)
	if err != nil {
		return nil, models.PaginationMeta{}, errors.New("post not found")
	}

	offset := (page - 1) * perPage
	comments, total, err := s.commentRepo.GetByPost(postID, offset, perPage)
	if err != nil {
		return nil, models.PaginationMeta{}, err
	}

	var responses []models.CommentResponse
	for _, comment := range comments {
		responses = append(responses, comment.ToResponse())
	}

	pagination := utils.CalculatePagination(page, perPage, total)
	return responses, pagination, nil
}

func (s *commentService) GetByAuthor(authorID uint, page, perPage int) ([]models.CommentResponse, models.PaginationMeta, error) {
	offset := (page - 1) * perPage
	comments, total, err := s.commentRepo.GetByAuthor(authorID, offset, perPage)
	if err != nil {
		return nil, models.PaginationMeta{}, err
	}

	var responses []models.CommentResponse
	for _, comment := range comments {
		responses = append(responses, comment.ToResponse())
	}

	pagination := utils.CalculatePagination(page, perPage, total)
	return responses, pagination, nil
}

func (s *commentService) GetPending(page, perPage int) ([]models.CommentResponse, models.PaginationMeta, error) {
	offset := (page - 1) * perPage
	comments, total, err := s.commentRepo.GetPending(offset, perPage)
	if err != nil {
		return nil, models.PaginationMeta{}, err
	}

	var responses []models.CommentResponse
	for _, comment := range comments {
		responses = append(responses, comment.ToResponse())
	}

	pagination := utils.CalculatePagination(page, perPage, total)
	return responses, pagination, nil
}

func (s *commentService) ApproveComment(commentID uint) (*models.CommentResponse, error) {
	_, err := s.commentRepo.GetByID(commentID)
	if err != nil {
		return nil, err
	}

	if err := s.commentRepo.UpdateStatus(commentID, models.CommentStatusApproved); err != nil {
		return nil, fmt.Errorf("failed to approve comment: %w", err)
	}

	// Get updated comment
	updatedComment, err := s.commentRepo.GetByID(commentID)
	if err != nil {
		return nil, err
	}

	response := updatedComment.ToResponse()
	return &response, nil
}

func (s *commentService) RejectComment(commentID uint) (*models.CommentResponse, error) {
	_, err := s.commentRepo.GetByID(commentID)
	if err != nil {
		return nil, err
	}

	if err := s.commentRepo.UpdateStatus(commentID, models.CommentStatusRejected); err != nil {
		return nil, fmt.Errorf("failed to reject comment: %w", err)
	}

	// Get updated comment
	updatedComment, err := s.commentRepo.GetByID(commentID)
	if err != nil {
		return nil, err
	}

	response := updatedComment.ToResponse()
	return &response, nil
}

func (s *commentService) GetPendingCount() (int64, error) {
	return s.commentRepo.CountPending()
}
