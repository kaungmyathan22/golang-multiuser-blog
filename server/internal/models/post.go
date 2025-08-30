package models

import (
	"time"
)

type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusPublished PostStatus = "published"
	PostStatusArchived  PostStatus = "archived"
)

type Post struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Title       string     `json:"title" gorm:"not null;size:200" validate:"required,min=5,max=200"`
	Slug        string     `json:"slug" gorm:"uniqueIndex;not null;size:250" validate:"required,min=5,max=250"`
	Content     string     `json:"content" gorm:"type:text;not null" validate:"required,min=10"`
	Excerpt     string     `json:"excerpt" gorm:"size:500" validate:"max=500"`
	FeaturedImg string     `json:"featured_image" gorm:"size:255" validate:"omitempty,url"`
	Status      PostStatus `json:"status" gorm:"default:'draft'" validate:"required,oneof=draft published archived"`
	ViewCount   int        `json:"view_count" gorm:"default:0"`
	AuthorID    uint       `json:"author_id" gorm:"not null" validate:"required"`
	PublishedAt *time.Time `json:"published_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Relationships
	Author   User      `json:"author" gorm:"foreignKey:AuthorID"`
	Comments []Comment `json:"comments,omitempty" gorm:"foreignKey:PostID"`
	Tags     []Tag     `json:"tags,omitempty" gorm:"many2many:post_tags;"`
}

// PostCreateRequest represents the request for creating a new post
type PostCreateRequest struct {
	Title       string     `json:"title" validate:"required,min=5,max=200"`
	Content     string     `json:"content" validate:"required,min=10"`
	Excerpt     string     `json:"excerpt" validate:"max=500"`
	FeaturedImg string     `json:"featured_image" validate:"omitempty,url"`
	Status      PostStatus `json:"status" validate:"required,oneof=draft published archived"`
	TagIDs      []uint     `json:"tag_ids" validate:"omitempty"`
}

// PostUpdateRequest represents the request for updating a post
type PostUpdateRequest struct {
	Title       string     `json:"title" validate:"omitempty,min=5,max=200"`
	Content     string     `json:"content" validate:"omitempty,min=10"`
	Excerpt     string     `json:"excerpt" validate:"max=500"`
	FeaturedImg string     `json:"featured_image" validate:"omitempty,url"`
	Status      PostStatus `json:"status" validate:"omitempty,oneof=draft published archived"`
	TagIDs      []uint     `json:"tag_ids" validate:"omitempty"`
}

// PostResponse represents the post response
type PostResponse struct {
	ID            uint          `json:"id"`
	Title         string        `json:"title"`
	Slug          string        `json:"slug"`
	Content       string        `json:"content"`
	Excerpt       string        `json:"excerpt"`
	FeaturedImg   string        `json:"featured_image"`
	Status        PostStatus    `json:"status"`
	ViewCount     int           `json:"view_count"`
	AuthorID      uint          `json:"author_id"`
	Author        UserResponse  `json:"author"`
	PublishedAt   *time.Time    `json:"published_at"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	Tags          []TagResponse `json:"tags,omitempty"`
	CommentsCount int           `json:"comments_count"`
}

// PostListResponse represents a simplified post response for listing
type PostListResponse struct {
	ID            uint          `json:"id"`
	Title         string        `json:"title"`
	Slug          string        `json:"slug"`
	Excerpt       string        `json:"excerpt"`
	FeaturedImg   string        `json:"featured_image"`
	Status        PostStatus    `json:"status"`
	ViewCount     int           `json:"view_count"`
	AuthorID      uint          `json:"author_id"`
	Author        UserResponse  `json:"author"`
	PublishedAt   *time.Time    `json:"published_at"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	Tags          []TagResponse `json:"tags,omitempty"`
	CommentsCount int           `json:"comments_count"`
}

// ToResponse converts Post to PostResponse
func (p *Post) ToResponse() PostResponse {
	return PostResponse{
		ID:          p.ID,
		Title:       p.Title,
		Slug:        p.Slug,
		Content:     p.Content,
		Excerpt:     p.Excerpt,
		FeaturedImg: p.FeaturedImg,
		Status:      p.Status,
		ViewCount:   p.ViewCount,
		AuthorID:    p.AuthorID,
		Author:      p.Author.ToResponse(),
		PublishedAt: p.PublishedAt,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

// ToListResponse converts Post to PostListResponse
func (p *Post) ToListResponse() PostListResponse {
	return PostListResponse{
		ID:          p.ID,
		Title:       p.Title,
		Slug:        p.Slug,
		Excerpt:     p.Excerpt,
		FeaturedImg: p.FeaturedImg,
		Status:      p.Status,
		ViewCount:   p.ViewCount,
		AuthorID:    p.AuthorID,
		Author:      p.Author.ToResponse(),
		PublishedAt: p.PublishedAt,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
