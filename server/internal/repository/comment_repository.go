package repository

import (
	"errors"

	"github.com/kaungmyathan22/golang-multiuser-blog/internal/models"
	"gorm.io/gorm"
)

type CommentRepository interface {
	Create(comment *models.Comment) error
	GetByID(id uint) (*models.Comment, error)
	Update(comment *models.Comment) error
	Delete(id uint) error
	GetByPost(postID uint, offset, limit int) ([]models.Comment, int64, error)
	GetByAuthor(authorID uint, offset, limit int) ([]models.Comment, int64, error)
	GetPending(offset, limit int) ([]models.Comment, int64, error)
	GetReplies(parentID uint) ([]models.Comment, error)
	CountByPost(postID uint) (int64, error)
	CountPending() (int64, error)
	UpdateStatus(id uint, status models.CommentStatus) error
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

func (r *commentRepository) GetByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	err := r.db.Preload("Author").Preload("Post").Preload("Replies", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Author").Where("status = ?", models.CommentStatusApproved)
	}).First(&comment, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("comment not found")
		}
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) Update(comment *models.Comment) error {
	return r.db.Save(comment).Error
}

func (r *commentRepository) Delete(id uint) error {
	// Delete all replies first
	if err := r.db.Where("parent_id = ?", id).Delete(&models.Comment{}).Error; err != nil {
		return err
	}

	// Delete the comment itself
	return r.db.Delete(&models.Comment{}, id).Error
}

func (r *commentRepository) GetByPost(postID uint, offset, limit int) ([]models.Comment, int64, error) {
	var comments []models.Comment
	var total int64

	query := r.db.Model(&models.Comment{}).Preload("Author").Preload("Replies", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Author").Where("status = ?", models.CommentStatusApproved).Order("created_at ASC")
	}).Where("post_id = ? AND parent_id IS NULL AND status = ?", postID, models.CommentStatusApproved)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&comments).Error
	return comments, total, err
}

func (r *commentRepository) GetByAuthor(authorID uint, offset, limit int) ([]models.Comment, int64, error) {
	var comments []models.Comment
	var total int64

	query := r.db.Model(&models.Comment{}).Preload("Author").Preload("Post").
		Where("author_id = ?", authorID)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&comments).Error
	return comments, total, err
}

func (r *commentRepository) GetPending(offset, limit int) ([]models.Comment, int64, error) {
	var comments []models.Comment
	var total int64

	query := r.db.Model(&models.Comment{}).Preload("Author").Preload("Post").
		Where("status = ?", models.CommentStatusPending)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Order("created_at ASC").Offset(offset).Limit(limit).Find(&comments).Error
	return comments, total, err
}

func (r *commentRepository) GetReplies(parentID uint) ([]models.Comment, error) {
	var replies []models.Comment
	err := r.db.Preload("Author").Where("parent_id = ? AND status = ?", parentID, models.CommentStatusApproved).
		Order("created_at ASC").Find(&replies).Error
	return replies, err
}

func (r *commentRepository) CountByPost(postID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Comment{}).Where("post_id = ? AND status = ?", postID, models.CommentStatusApproved).Count(&count).Error
	return count, err
}

func (r *commentRepository) CountPending() (int64, error) {
	var count int64
	err := r.db.Model(&models.Comment{}).Where("status = ?", models.CommentStatusPending).Count(&count).Error
	return count, err
}

func (r *commentRepository) UpdateStatus(id uint, status models.CommentStatus) error {
	return r.db.Model(&models.Comment{}).Where("id = ?", id).Update("status", status).Error
}
