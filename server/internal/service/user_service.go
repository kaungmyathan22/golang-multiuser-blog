package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/kaungmyathan22/golang-multiuser-blog/internal/config"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/models"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/repository"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/utils"
)

type UserService interface {
	Register(req *models.UserCreateRequest) (*models.UserResponse, error)
	Login(req *models.UserLoginRequest) (*models.AuthResponse, error)
	GetProfile(userID uint) (*models.UserResponse, error)
	UpdateProfile(userID uint, req *models.UserUpdateRequest) (*models.UserResponse, error)
	GetUsers(page, perPage int) ([]models.UserResponse, models.PaginationMeta, error)
	GetUserByID(id uint) (*models.UserResponse, error)
	DeactivateUser(id uint) error
	ActivateUser(id uint) error
	ChangePassword(userID uint, oldPassword, newPassword string) error
	RefreshToken(token string) (*models.AuthResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
	config   *config.Config
}

func NewUserService(userRepo repository.UserRepository, config *config.Config) UserService {
	return &userService{
		userRepo: userRepo,
		config:   config,
	}
}

func (s *userService) Register(req *models.UserCreateRequest) (*models.UserResponse, error) {
	// Validate request
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, fmt.Errorf("validation failed: %v", validationErrors)
	}

	// Check if email is already taken
	if s.userRepo.IsEmailTaken(req.Email, 0) {
		return nil, errors.New("email is already registered")
	}

	// Check if username is already taken
	if s.userRepo.IsUsernameTaken(req.Username, 0) {
		return nil, errors.New("username is already taken")
	}

	// Create user
	user := &models.User{
		FirstName: utils.SanitizeText(req.FirstName),
		LastName:  utils.SanitizeText(req.LastName),
		Email:     req.Email,
		Username:  req.Username,
		Password:  req.Password, // Will be hashed by BeforeCreate hook
		Bio:       utils.SanitizeText(req.Bio),
		Avatar:    req.Avatar,
		IsActive:  true,
		IsAdmin:   false,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *userService) Login(req *models.UserLoginRequest) (*models.AuthResponse, error) {
	// Validate request
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, fmt.Errorf("validation failed: %v", validationErrors)
	}

	// Find user by email or username
	user, err := s.userRepo.GetByEmailOrUsername(req.EmailOrUsername)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Verify password
	if !user.CheckPassword(req.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user, s.config)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	userResponse := user.ToResponse()
	return &models.AuthResponse{
		User:      userResponse,
		Token:     token,
		TokenType: "Bearer",
		ExpiresIn: int(s.config.JWT.ExpiresIn.Seconds()),
	}, nil
}

func (s *userService) GetProfile(userID uint) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *userService) UpdateProfile(userID uint, req *models.UserUpdateRequest) (*models.UserResponse, error) {
	// Validate request
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, fmt.Errorf("validation failed: %v", validationErrors)
	}

	// Get existing user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Check if email is already taken (excluding current user)
	if req.Email != "" && req.Email != user.Email {
		if s.userRepo.IsEmailTaken(req.Email, userID) {
			return nil, errors.New("email is already registered")
		}
		user.Email = req.Email
	}

	// Check if username is already taken (excluding current user)
	if req.Username != "" && req.Username != user.Username {
		if s.userRepo.IsUsernameTaken(req.Username, userID) {
			return nil, errors.New("username is already taken")
		}
		user.Username = req.Username
	}

	// Update other fields
	if req.FirstName != "" {
		user.FirstName = utils.SanitizeText(req.FirstName)
	}
	if req.LastName != "" {
		user.LastName = utils.SanitizeText(req.LastName)
	}
	user.Bio = utils.SanitizeText(req.Bio)
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *userService) GetUsers(page, perPage int) ([]models.UserResponse, models.PaginationMeta, error) {
	offset := (page - 1) * perPage
	users, total, err := s.userRepo.List(offset, perPage)
	if err != nil {
		return nil, models.PaginationMeta{}, err
	}

	var responses []models.UserResponse
	for _, user := range users {
		responses = append(responses, user.ToResponse())
	}

	pagination := utils.CalculatePagination(page, perPage, total)
	return responses, pagination, nil
}

func (s *userService) GetUserByID(id uint) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *userService) DeactivateUser(id uint) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	user.IsActive = false
	return s.userRepo.Update(user)
}

func (s *userService) ActivateUser(id uint) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	user.IsActive = true
	return s.userRepo.Update(user)
}

func (s *userService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// Verify old password
	if !user.CheckPassword(oldPassword) {
		return errors.New("invalid current password")
	}

	// Validate new password
	if len(newPassword) < 8 {
		return errors.New("new password must be at least 8 characters long")
	}

	// Update password (will be hashed by BeforeCreate hook)
	user.Password = newPassword
	user.UpdatedAt = time.Now()

	return s.userRepo.Update(user)
}

func (s *userService) RefreshToken(token string) (*models.AuthResponse, error) {
	// Validate and refresh token
	newToken, err := utils.RefreshToken(token, s.config)
	if err != nil {
		return nil, err
	}

	// Extract user info from token
	claims, err := utils.ValidateToken(newToken, s.config)
	if err != nil {
		return nil, err
	}

	// Get fresh user data
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	// Check if user is still active
	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	userResponse := user.ToResponse()
	return &models.AuthResponse{
		User:      userResponse,
		Token:     newToken,
		TokenType: "Bearer",
		ExpiresIn: int(s.config.JWT.ExpiresIn.Seconds()),
	}, nil
}
