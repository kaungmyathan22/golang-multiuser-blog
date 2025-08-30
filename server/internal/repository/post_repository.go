package repository

import (
	"errors"
	"strings"
	"time"

	"github.com/kaungmyathan22/golang-multiuser-blog/internal/models"
	"gorm.io/gorm"
)

type PostRepository interface {
	Create(post *models.Post) error
	GetByID(id uint) (*models.Post, error)
	GetBySlug(slug string) (*models.Post, error)
	Update(post *models.Post) error
	Delete(id uint) error
	List(offset, limit int, status models.PostStatus, authorID uint) ([]models.Post, int64, error)
	GetPublished(offset, limit int) ([]models.Post, int64, error)
	GetByAuthor(authorID uint, offset, limit int) ([]models.Post, int64, error)
	GetByTag(tagID uint, offset, limit int) ([]models.Post, int64, error)
	Search(query string, offset, limit int) ([]models.Post, int64, error)
	IncrementViewCount(id uint) error
	IsSlugTaken(slug string, excludeID uint) bool
	AddTags(postID uint, tagIDs []uint) error
	RemoveTags(postID uint, tagIDs []uint) error
	UpdateTags(postID uint, tagIDs []uint) error
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) GetByID(id uint) (*models.Post, error) {
	var post models.Post
	err := r.db.Preload("Author").Preload("Tags").Preload("Comments", func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", models.CommentStatusApproved)
	}).First(&post, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) GetBySlug(slug string) (*models.Post, error) {
	var post models.Post
	err := r.db.Preload("Author").Preload("Tags").Preload("Comments", func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", models.CommentStatusApproved)
	}).Where("slug = ?", slug).First(&post).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) Update(post *models.Post) error {
	return r.db.Save(post).Error
}

func (r *postRepository) Delete(id uint) error {
	return r.db.Delete(&models.Post{}, id).Error
}

func (r *postRepository) List(offset, limit int, status models.PostStatus, authorID uint) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	query := r.db.Model(&models.Post{}).Preload("Author").Preload("Tags")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if authorID > 0 {
		query = query.Where("author_id = ?", authorID)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&posts).Error
	return posts, total, err
}

func (r *postRepository) GetPublished(offset, limit int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	query := r.db.Model(&models.Post{}).Preload("Author").Preload("Tags").
		Where("status = ? AND published_at <= ?", models.PostStatusPublished, time.Now())

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Order("published_at DESC").Offset(offset).Limit(limit).Find(&posts).Error
	return posts, total, err
}

func (r *postRepository) GetByAuthor(authorID uint, offset, limit int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	query := r.db.Model(&models.Post{}).Preload("Author").Preload("Tags").
		Where("author_id = ? AND status = ?", authorID, models.PostStatusPublished)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Order("published_at DESC").Offset(offset).Limit(limit).Find(&posts).Error
	return posts, total, err
}

func (r *postRepository) GetByTag(tagID uint, offset, limit int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	subQuery := r.db.Table("post_tags").Select("post_id").Where("tag_id = ?", tagID)
	query := r.db.Model(&models.Post{}).Preload("Author").Preload("Tags").
		Where("id IN (?) AND status = ?", subQuery, models.PostStatusPublished)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Order("published_at DESC").Offset(offset).Limit(limit).Find(&posts).Error
	return posts, total, err
}

func (r *postRepository) Search(query string, offset, limit int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	searchQuery := "%" + strings.ToLower(query) + "%"
	dbQuery := r.db.Model(&models.Post{}).Preload("Author").Preload("Tags").
		Where("status = ? AND (LOWER(title) LIKE ? OR LOWER(content) LIKE ? OR LOWER(excerpt) LIKE ?)",
			models.PostStatusPublished, searchQuery, searchQuery, searchQuery)

	// Count total records
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := dbQuery.Order("published_at DESC").Offset(offset).Limit(limit).Find(&posts).Error
	return posts, total, err
}

func (r *postRepository) IncrementViewCount(id uint) error {
	return r.db.Model(&models.Post{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

func (r *postRepository) IsSlugTaken(slug string, excludeID uint) bool {
	var count int64
	query := r.db.Model(&models.Post{}).Where("slug = ?", slug)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	query.Count(&count)
	return count > 0
}

func (r *postRepository) AddTags(postID uint, tagIDs []uint) error {
	var post models.Post
	if err := r.db.First(&post, postID).Error; err != nil {
		return err
	}

	var tags []models.Tag
	if err := r.db.Find(&tags, tagIDs).Error; err != nil {
		return err
	}

	return r.db.Model(&post).Association("Tags").Append(&tags)
}

func (r *postRepository) RemoveTags(postID uint, tagIDs []uint) error {
	var post models.Post
	if err := r.db.First(&post, postID).Error; err != nil {
		return err
	}

	var tags []models.Tag
	if err := r.db.Find(&tags, tagIDs).Error; err != nil {
		return err
	}

	return r.db.Model(&post).Association("Tags").Delete(&tags)
}

func (r *postRepository) UpdateTags(postID uint, tagIDs []uint) error {
	var post models.Post
	if err := r.db.First(&post, postID).Error; err != nil {
		return err
	}

	var tags []models.Tag
	if err := r.db.Find(&tags, tagIDs).Error; err != nil {
		return err
	}

	return r.db.Model(&post).Association("Tags").Replace(&tags)
}
