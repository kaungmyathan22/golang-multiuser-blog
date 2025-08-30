package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	FirstName string    `json:"first_name" gorm:"not null;size:50" validate:"required,min=2,max=50"`
	LastName  string    `json:"last_name" gorm:"not null;size:50" validate:"required,min=2,max=50"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null;size:100" validate:"required,email,max=100"`
	Username  string    `json:"username" gorm:"uniqueIndex;not null;size:30" validate:"required,min=3,max=30,alphanum"`
	Password  string    `json:"-" gorm:"not null" validate:"required,min=8"`
	Bio       string    `json:"bio" gorm:"size:500" validate:"max=500"`
	Avatar    string    `json:"avatar" gorm:"size:255" validate:"omitempty,url"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	IsAdmin   bool      `json:"is_admin" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	Posts    []Post    `json:"posts,omitempty" gorm:"foreignKey:AuthorID"`
	Comments []Comment `json:"comments,omitempty" gorm:"foreignKey:AuthorID"`
}

// UserCreateRequest represents the request for creating a new user
type UserCreateRequest struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name" validate:"required,min=2,max=50"`
	Email     string `json:"email" validate:"required,email,max=100"`
	Username  string `json:"username" validate:"required,min=3,max=30,alphanum"`
	Password  string `json:"password" validate:"required,min=8"`
	Bio       string `json:"bio" validate:"max=500"`
	Avatar    string `json:"avatar" validate:"omitempty,url"`
}

// UserUpdateRequest represents the request for updating user data
type UserUpdateRequest struct {
	FirstName string `json:"first_name" validate:"omitempty,min=2,max=50"`
	LastName  string `json:"last_name" validate:"omitempty,min=2,max=50"`
	Email     string `json:"email" validate:"omitempty,email,max=100"`
	Username  string `json:"username" validate:"omitempty,min=3,max=30,alphanum"`
	Bio       string `json:"bio" validate:"max=500"`
	Avatar    string `json:"avatar" validate:"omitempty,url"`
}

// UserLoginRequest represents the login request
type UserLoginRequest struct {
	EmailOrUsername string `json:"email_or_username" validate:"required"`
	Password        string `json:"password" validate:"required"`
}

// UserResponse represents the user response (without sensitive data)
type UserResponse struct {
	ID        uint      `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Bio       string    `json:"bio"`
	Avatar    string    `json:"avatar"`
	IsActive  bool      `json:"is_active"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate is a GORM hook that runs before creating a user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// CheckPassword verifies if the provided password matches the hashed password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Username:  u.Username,
		Bio:       u.Bio,
		Avatar:    u.Avatar,
		IsActive:  u.IsActive,
		IsAdmin:   u.IsAdmin,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
