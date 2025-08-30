package models

import (
	"time"
)

type Tag struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"uniqueIndex;not null;size:50" validate:"required,min=2,max=50"`
	Slug        string    `json:"slug" gorm:"uniqueIndex;not null;size:60" validate:"required,min=2,max=60"`
	Description string    `json:"description" gorm:"size:200" validate:"max=200"`
	Color       string    `json:"color" gorm:"size:7;default:'#3B82F6'" validate:"omitempty,hexcolor"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Posts []Post `json:"posts,omitempty" gorm:"many2many:post_tags;"`
}

// TagCreateRequest represents the request for creating a new tag
type TagCreateRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=50"`
	Description string `json:"description" validate:"max=200"`
	Color       string `json:"color" validate:"omitempty,hexcolor"`
}

// TagUpdateRequest represents the request for updating a tag
type TagUpdateRequest struct {
	Name        string `json:"name" validate:"omitempty,min=2,max=50"`
	Description string `json:"description" validate:"max=200"`
	Color       string `json:"color" validate:"omitempty,hexcolor"`
}

// TagResponse represents the tag response
type TagResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
	PostsCount  int       `json:"posts_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToResponse converts Tag to TagResponse
func (t *Tag) ToResponse() TagResponse {
	return TagResponse{
		ID:          t.ID,
		Name:        t.Name,
		Slug:        t.Slug,
		Description: t.Description,
		Color:       t.Color,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}