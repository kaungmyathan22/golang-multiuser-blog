package handlers
package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/handlers"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/models"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockUserService is a mock implementation of the UserService interface
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Register(req *models.UserCreateRequest) (*models.UserResponse, error) {
	args := m.Called(req)
	return args.Get(0).(*models.UserResponse), args.Error(1)
}

func (m *MockUserService) Login(req *models.UserLoginRequest) (*models.AuthResponse, error) {
	args := m.Called(req)
	return args.Get(0).(*models.AuthResponse), args.Error(1)
}

func (m *MockUserService) GetProfile(userID uint) (*models.UserResponse, error) {
	args := m.Called(userID)
	return args.Get(0).(*models.UserResponse), args.Error(1)
}

func (m *MockUserService) UpdateProfile(userID uint, req *models.UserUpdateRequest) (*models.UserResponse, error) {
	args := m.Called(userID, req)
	return args.Get(0).(*models.UserResponse), args.Error(1)
}

func (m *MockUserService) GetUsers(page, perPage int) ([]models.UserResponse, models.PaginationMeta, error) {
	args := m.Called(page, perPage)
	return args.Get(0).([]models.UserResponse), args.Get(1).(models.PaginationMeta), args.Error(2)
}

func (m *MockUserService) GetUserByID(id uint) (*models.UserResponse, error) {
	args := m.Called(id)
	return args.Get(0).(*models.UserResponse), args.Error(1)
}

func (m *MockUserService) DeactivateUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserService) ActivateUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	args := m.Called(userID, oldPassword, newPassword)
	return args.Error(0)
}

func (m *MockUserService) RefreshToken(token string) (*models.AuthResponse, error) {
	args := m.Called(token)
	return args.Get(0).(*models.AuthResponse), args.Error(1)
}

func TestAuthHandler_Register(t *testing.T) {
	// Skip integration tests in short mode
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	gin.SetMode(gin.TestMode)

	t.Run("successful registration", func(t *testing.T) {
		// Create mock service
		mockService := new(MockUserService)
		handler := handlers.NewAuthHandler(mockService)

		// Create test request
		userReq := &models.UserCreateRequest{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Username:  "johndoe",
			Password:  "password123",
		}

		// Set up mock expectations
		mockService.On("Register", mock.AnythingOfType("*models.UserCreateRequest")).Return(&models.UserResponse{
			ID:        1,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Username:  "johndoe",
		}, nil)

		// Create HTTP request
		jsonReq, _ := json.Marshal(userReq)
		req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonReq))
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		w := httptest.NewRecorder()

		// Create gin context and serve HTTP
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call handler
		handler.Register(c)

		// Assertions
		require.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid request format", func(t *testing.T) {
		// Create mock service
		mockService := new(MockUserService)
		handler := handlers.NewAuthHandler(mockService)

		// Create invalid JSON request
		invalidJSON := []byte(`{"invalid": json}`)
		req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(invalidJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		w := httptest.NewRecorder()

		// Create gin context and serve HTTP
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call handler
		handler.Register(c)

		// Assertions
		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAuthHandler_Login(t *testing.T) {
	// Skip integration tests in short mode
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	gin.SetMode(gin.TestMode)

	t.Run("successful login", func(t *testing.T) {
		// Create mock service
		mockService := new(MockUserService)
		handler := handlers.NewAuthHandler(mockService)

		// Create test request
		loginReq := &models.UserLoginRequest{
			EmailOrUsername: "johndoe",
			Password:        "password123",
		}

		// Set up mock expectations
		mockService.On("Login", mock.AnythingOfType("*models.UserLoginRequest")).Return(&models.AuthResponse{
			User: models.UserResponse{
				ID:        1,
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				Username:  "johndoe",
			},
			Token:     "test.jwt.token",
			TokenType: "Bearer",
			ExpiresIn: 86400,
		}, nil)

		// Create HTTP request
		jsonReq, _ := json.Marshal(loginReq)
		req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonReq))
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		w := httptest.NewRecorder()

		// Create gin context and serve HTTP
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call handler
		handler.Login(c)

		// Assertions
		require.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})
}