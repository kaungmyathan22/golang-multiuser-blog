package e2e
//go:build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/kaungmyathan22/golang-multiuser-blog/internal/config"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/migration"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/models"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/router"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testServer *httptest.Server
)

func TestMain(m *testing.M) {
	// Set test environment
	os.Setenv("GIN_MODE", "test")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_NAME", "golang_multiuser_blog_e2e_test")
	os.Setenv("DB_SSLMODE", "disable")
	os.Setenv("JWT_SECRET", "test-secret-key-for-e2e")
	os.Setenv("PORT", "8081")

	// Load test config
	cfg := config.LoadConfig()

	// Initialize database
	config.InitDatabase(cfg)

	// Run migrations
	err := migration.RunMigrations()
	if err != nil {
		panic("Failed to run migrations: " + err.Error())
	}

	// Initialize router
	r := router.NewRouter(cfg)
	appRouter := r.SetupRoutes()

	// Create test server
	testServer = httptest.NewServer(appRouter)

	// Run tests
	code := m.Run()

	// Clean up
	testServer.Close()

	os.Exit(code)
}

func TestAuthFlow(t *testing.T) {
	// Test user registration
	t.Run("RegisterUser", func(t *testing.T) {
		userReq := &models.UserCreateRequest{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe.e2e@example.com",
			Username:  "johndoe_e2e",
			Password:  "password123",
			Bio:       "Test user for E2E testing",
		}

		jsonReq, _ := json.Marshal(userReq)
		resp, err := http.Post(testServer.URL+"/api/auth/register", "application/json", bytes.NewBuffer(jsonReq))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var registerResp models.APIResponse
		err = json.NewDecoder(resp.Body).Decode(&registerResp)
		require.NoError(t, err)
		assert.True(t, registerResp.Success)
	})

	// Test user login
	t.Run("LoginUser", func(t *testing.T) {
		loginReq := &models.UserLoginRequest{
			EmailOrUsername: "johndoe_e2e",
			Password:        "password123",
		}

		jsonReq, _ := json.Marshal(loginReq)
		resp, err := http.Post(testServer.URL+"/api/auth/login", "application/json", bytes.NewBuffer(jsonReq))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var loginResp models.APIResponse
		err = json.NewDecoder(resp.Body).Decode(&loginResp)
		require.NoError(t, err)
		assert.True(t, loginResp.Success)

		// Extract token for future requests
		authData, ok := loginResp.Data.(map[string]interface{})
		require.True(t, ok)
		token, ok := authData["token"].(string)
		require.True(t, ok)
		assert.NotEmpty(t, token)

		// Test get profile with token
		t.Run("GetProfile", func(t *testing.T) {
			req, _ := http.NewRequest("GET", testServer.URL+"/api/auth/profile", nil)
			req.Header.Set("Authorization", "Bearer "+token)

			client := &http.Client{Timeout: 10 * time.Second}
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var profileResp models.APIResponse
			err = json.NewDecoder(resp.Body).Decode(&profileResp)
			require.NoError(t, err)
			assert.True(t, profileResp.Success)
		})
	})
}

func TestHealthCheck(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/health")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var healthResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&healthResp)
	require.NoError(t, err)

	success, ok := healthResp["success"].(bool)
	require.True(t, ok)
	assert.True(t, success)
}