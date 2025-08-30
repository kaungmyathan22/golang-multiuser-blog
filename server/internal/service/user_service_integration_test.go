//go:build integration

package service_test

import (
	"os"
	"testing"

	"github.com/kaungmyathan22/golang-multiuser-blog/internal/config"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/models"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/repository"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var (
	testDB   *gorm.DB
	userRepo repository.UserRepository
	userSvc  service.UserService
)

func TestMain(m *testing.M) {
	// Set test environment
	os.Setenv("GIN_MODE", "test")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_NAME", "golang_multiuser_blog_test")
	os.Setenv("DB_SSLMODE", "disable")
	os.Setenv("JWT_SECRET", "test-secret-key")

	// Load test config
	cfg := config.LoadConfig()

	// Initialize test database
	config.InitDatabase(cfg)
	testDB = config.GetDB()

	// Run migrations
	err := testDB.AutoMigrate(&models.User{})
	if err != nil {
		panic("Failed to run migrations: " + err.Error())
	}

	// Initialize repositories and services
	userRepo = repository.NewUserRepository(testDB)
	userSvc = service.NewUserService(userRepo, cfg)

	// Run tests
	code := m.Run()

	// Clean up
	testDB.Migrator().DropTable(&models.User{})

	os.Exit(code)
}

func TestUserService_Register_Success(t *testing.T) {
	// Prepare test data
	req := &models.UserCreateRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Username:  "johndoe",
		Password:  "password123",
		Bio:       "Test user bio",
	}

	// Execute
	user, err := userSvc.Register(req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, req.FirstName, user.FirstName)
	assert.Equal(t, req.LastName, user.LastName)
	assert.Equal(t, req.Email, user.Email)
	assert.Equal(t, req.Username, user.Username)
	assert.Empty(t, user.Password) // Password should not be returned
}

func TestUserService_Register_DuplicateEmail(t *testing.T) {
	// Prepare test data
	req1 := &models.UserCreateRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe2@example.com",
		Username:  "johndoe2",
		Password:  "password123",
	}

	req2 := &models.UserCreateRequest{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "john.doe2@example.com", // Same email
		Username:  "janedoe",
		Password:  "password123",
	}

	// Register first user
	_, err1 := userSvc.Register(req1)
	require.NoError(t, err1)

	// Try to register second user with same email
	_, err2 := userSvc.Register(req2)

	// Assert
	assert.Error(t, err2)
	assert.Contains(t, err2.Error(), "email is already registered")
}

func TestUserService_Login_Success(t *testing.T) {
	// Prepare test data
	registerReq := &models.UserCreateRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe3@example.com",
		Username:  "johndoe3",
		Password:  "password123",
	}

	loginReq := &models.UserLoginRequest{
		EmailOrUsername: "johndoe3",
		Password:        "password123",
	}

	// Register user
	_, err := userSvc.Register(registerReq)
	require.NoError(t, err)

	// Login
	authResponse, err := userSvc.Login(loginReq)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, authResponse)
	assert.NotEmpty(t, authResponse.Token)
	assert.Equal(t, "Bearer", authResponse.TokenType)
	assert.Greater(t, authResponse.ExpiresIn, 0)
}

func TestUserService_Login_InvalidCredentials(t *testing.T) {
	// Prepare test data
	loginReq := &models.UserLoginRequest{
		EmailOrUsername: "nonexistentuser",
		Password:        "wrongpassword",
	}

	// Login with invalid credentials
	_, err := userSvc.Login(loginReq)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
}
