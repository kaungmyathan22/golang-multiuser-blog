package models

import (
	"time"
)

type CommentStatus string

const (
	CommentStatusPending  CommentStatus = "pending"
	CommentStatusApproved CommentStatus = "approved"
	CommentStatusRejected CommentStatus = "rejected"
)

type Comment struct {
	ID        uint          `json:"id" gorm:"primaryKey"`
	Content   string        `json:"content" gorm:"type:text;not null" validate:"required,min=1,max=1000"`
	Status    CommentStatus `json:"status" gorm:"default:'pending'" validate:"oneof=pending approved rejected"`
	AuthorID  uint          `json:"author_id" gorm:"not null" validate:"required"`
	PostID    uint          `json:"post_id" gorm:"not null" validate:"required"`
	ParentID  *uint         `json:"parent_id" gorm:"index"` // For nested comments/replies
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`

	// Relationships
	Author  User      `json:"author" gorm:"foreignKey:AuthorID"`
	Post    Post      `json:"post" gorm:"foreignKey:PostID"`
	Parent  *Comment  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Replies []Comment `json:"replies,omitempty" gorm:"foreignKey:ParentID"`
}

// CommentCreateRequest represents the request for creating a new comment
type CommentCreateRequest struct {
	Content  string `json:"content" validate:"required,min=1,max=1000"`
	PostID   uint   `json:"post_id" validate:"required"`
	ParentID *uint  `json:"parent_id" validate:"omitempty"`
}

// CommentUpdateRequest represents the request for updating a comment
type CommentUpdateRequest struct {
	Content string        `json:"content" validate:"omitempty,min=1,max=1000"`
	Status  CommentStatus `json:"status" validate:"omitempty,oneof=pending approved rejected"`
}

// CommentResponse represents the comment response
type CommentResponse struct {
	ID        uint              `json:"id"`
	Content   string            `json:"content"`
	Status    CommentStatus     `json:"status"`
	AuthorID  uint              `json:"author_id"`
	PostID    uint              `json:"post_id"`
	ParentID  *uint             `json:"parent_id"`
	Author    UserResponse      `json:"author"`
	Replies   []CommentResponse `json:"replies,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// ToResponse converts Comment to CommentResponse
func (c *Comment) ToResponse() CommentResponse {
	response := CommentResponse{
		ID:        c.ID,
		Content:   c.Content,
		Status:    c.Status,
		AuthorID:  c.AuthorID,
		PostID:    c.PostID,
		ParentID:  c.ParentID,
		Author:    c.Author.ToResponse(),
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}

	// Convert replies if they exist
	if len(c.Replies) > 0 {
		response.Replies = make([]CommentResponse, len(c.Replies))
		for i, reply := range c.Replies {
			response.Replies[i] = reply.ToResponse()
		}
	}

	return response
}
