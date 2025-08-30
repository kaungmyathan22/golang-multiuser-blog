package repository

import (
	"errors"

	"github.com/kaungmyathan22/golang-multiuser-blog/internal/models"
	"gorm.io/gorm"
)

type TagRepository interface {
	Create(tag *models.Tag) error
	GetByID(id uint) (*models.Tag, error)
	GetBySlug(slug string) (*models.Tag, error)
	Update(tag *models.Tag) error
	Delete(id uint) error
	List(offset, limit int) ([]models.Tag, int64, error)
	GetAll() ([]models.Tag, error)
	IsNameTaken(name string, excludeID uint) bool
	IsSlugTaken(slug string, excludeID uint) bool
	GetPopular(limit int) ([]models.Tag, error)
}

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) Create(tag *models.Tag) error {
	return r.db.Create(tag).Error
}

func (r *tagRepository) GetByID(id uint) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.Preload("Posts", func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", models.PostStatusPublished)
	}).First(&tag, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tag not found")
		}
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) GetBySlug(slug string) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.Preload("Posts", func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", models.PostStatusPublished)
	}).Where("slug = ?", slug).First(&tag).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tag not found")
		}
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) Update(tag *models.Tag) error {
	return r.db.Save(tag).Error
}

func (r *tagRepository) Delete(id uint) error {
	// Remove associations with posts first
	var tag models.Tag
	if err := r.db.First(&tag, id).Error; err != nil {
		return err
	}

	if err := r.db.Model(&tag).Association("Posts").Clear(); err != nil {
		return err
	}

	// Delete the tag
	return r.db.Delete(&models.Tag{}, id).Error
}

func (r *tagRepository) List(offset, limit int) ([]models.Tag, int64, error) {
	var tags []models.Tag
	var total int64

	// Count total records
	if err := r.db.Model(&models.Tag{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results with post count
	err := r.db.Select("tags.*, COUNT(post_tags.post_id) as posts_count").
		Joins("LEFT JOIN post_tags ON tags.id = post_tags.tag_id").
		Group("tags.id").
		Order("tags.name ASC").
		Offset(offset).
		Limit(limit).
		Find(&tags).Error

	return tags, total, err
}

func (r *tagRepository) GetAll() ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.Order("name ASC").Find(&tags).Error
	return tags, err
}

func (r *tagRepository) IsNameTaken(name string, excludeID uint) bool {
	var count int64
	query := r.db.Model(&models.Tag{}).Where("LOWER(name) = LOWER(?)", name)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	query.Count(&count)
	return count > 0
}

func (r *tagRepository) IsSlugTaken(slug string, excludeID uint) bool {
	var count int64
	query := r.db.Model(&models.Tag{}).Where("slug = ?", slug)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	query.Count(&count)
	return count > 0
}

func (r *tagRepository) GetPopular(limit int) ([]models.Tag, error) {
	var tags []models.Tag

	err := r.db.Select("tags.*, COUNT(post_tags.post_id) as posts_count").
		Joins("LEFT JOIN post_tags ON tags.id = post_tags.tag_id").
		Joins("LEFT JOIN posts ON post_tags.post_id = posts.id AND posts.status = ?", models.PostStatusPublished).
		Group("tags.id").
		Having("COUNT(post_tags.post_id) > 0").
		Order("posts_count DESC").
		Limit(limit).
		Find(&tags).Error

	return tags, err
}
