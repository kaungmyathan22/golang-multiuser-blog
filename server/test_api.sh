#!/bin/bash

# Golang Multi-User Blog API Test Script
# This script tests the main API endpoints

BASE_URL="http://localhost:8080"
API_URL="$BASE_URL/api"

echo "🧪 Testing Golang Multi-User Blog API"
echo "======================================"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ PASS${NC}: $2"
    else
        echo -e "${RED}❌ FAIL${NC}: $2"
    fi
}

print_info() {
    echo -e "${BLUE}ℹ️  INFO${NC}: $1"
}

print_warning() {
    echo -e "${YELLOW}⚠️  WARN${NC}: $1"
}

# Check if server is running
echo "1. Testing server health..."
HEALTH_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/health")
if [ "$HEALTH_RESPONSE" = "200" ]; then
    print_status 0 "Server is healthy"
else
    print_status 1 "Server is not responding (HTTP $HEALTH_RESPONSE)"
    echo "Please make sure the server is running on $BASE_URL"
    exit 1
fi

echo ""

# Test user registration
echo "2. Testing user registration..."
REGISTER_RESPONSE=$(curl -s -X POST "$API_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d '{
        "first_name": "Test",
        "last_name": "User",
        "email": "test@example.com",
        "username": "testuser",
        "password": "testpass123",
        "bio": "Test user for API testing"
    }' \
    -w "%{http_code}")

HTTP_CODE=$(echo "$REGISTER_RESPONSE" | tail -c 4)
if [ "$HTTP_CODE" = "201" ] || [ "$HTTP_CODE" = "409" ]; then
    if [ "$HTTP_CODE" = "201" ]; then
        print_status 0 "User registration successful"
    else
        print_warning "User already exists (expected if running multiple times)"
    fi
else
    print_status 1 "User registration failed (HTTP $HTTP_CODE)"
fi

echo ""

# Test user login
echo "3. Testing user login..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d '{
        "email_or_username": "test@example.com",
        "password": "testpass123"
    }')

# Extract token from response
TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ ! -z "$TOKEN" ]; then
    print_status 0 "User login successful"
    print_info "JWT Token obtained: ${TOKEN:0:20}..."
else
    print_status 1 "User login failed"
    echo "Response: $LOGIN_RESPONSE"
fi

echo ""

# Test getting user profile
echo "4. Testing user profile..."
if [ ! -z "$TOKEN" ]; then
    PROFILE_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "$API_URL/auth/profile" \
        -H "Authorization: Bearer $TOKEN")

    if [ "$PROFILE_RESPONSE" = "200" ]; then
        print_status 0 "Profile retrieval successful"
    else
        print_status 1 "Profile retrieval failed (HTTP $PROFILE_RESPONSE)"
    fi
else
    print_warning "Skipping profile test (no token available)"
fi

echo ""

# Test creating a post
echo "5. Testing post creation..."
if [ ! -z "$TOKEN" ]; then
    POST_RESPONSE=$(curl -s -X POST "$API_URL/posts" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d '{
            "title": "Test Blog Post",
            "content": "This is a test blog post created by the API test script. It contains some sample content to validate the post creation functionality.",
            "excerpt": "A test blog post for API validation",
            "status": "published"
        }' \
        -w "%{http_code}")

    HTTP_CODE=$(echo "$POST_RESPONSE" | tail -c 4)
    if [ "$HTTP_CODE" = "201" ]; then
        print_status 0 "Post creation successful"
        # Extract post ID for later tests
        POST_ID=$(echo "$POST_RESPONSE" | grep -o '"id":[0-9]*' | cut -d':' -f2 | head -1)
        print_info "Created post ID: $POST_ID"
    else
        print_status 1 "Post creation failed (HTTP $HTTP_CODE)"
    fi
else
    print_warning "Skipping post creation test (no token available)"
fi

echo ""

# Test getting published posts
echo "6. Testing published posts retrieval..."
POSTS_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "$API_URL/posts/published?page=1&per_page=5")

if [ "$POSTS_RESPONSE" = "200" ]; then
    print_status 0 "Published posts retrieval successful"
else
    print_status 1 "Published posts retrieval failed (HTTP $POSTS_RESPONSE)"
fi

echo ""

# Test post search
echo "7. Testing post search..."
SEARCH_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "$API_URL/posts/search?q=test&page=1&per_page=5")

if [ "$SEARCH_RESPONSE" = "200" ]; then
    print_status 0 "Post search successful"
else
    print_status 1 "Post search failed (HTTP $SEARCH_RESPONSE)"
fi

echo ""

# Test admin login
echo "8. Testing admin login..."
ADMIN_LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d '{
        "email_or_username": "admin@blog.com",
        "password": "admin123456"
    }')

ADMIN_TOKEN=$(echo "$ADMIN_LOGIN_RESPONSE" | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ ! -z "$ADMIN_TOKEN" ]; then
    print_status 0 "Admin login successful"
    print_info "Admin JWT Token obtained: ${ADMIN_TOKEN:0:20}..."
else
    print_status 1 "Admin login failed"
    print_warning "Make sure the server has been started and migrations have run"
fi

echo ""

# Test unauthorized access
echo "9. Testing unauthorized access protection..."
UNAUTH_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "$API_URL/auth/profile")

if [ "$UNAUTH_RESPONSE" = "401" ]; then
    print_status 0 "Unauthorized access properly blocked"
else
    print_status 1 "Unauthorized access not properly blocked (HTTP $UNAUTH_RESPONSE)"
fi

echo ""

# Test invalid endpoint
echo "10. Testing invalid endpoint..."
INVALID_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "$API_URL/invalid-endpoint")

if [ "$INVALID_RESPONSE" = "404" ]; then
    print_status 0 "Invalid endpoint properly returns 404"
else
    print_status 1 "Invalid endpoint handling unexpected (HTTP $INVALID_RESPONSE)"
fi

echo ""
echo "🏁 API Testing Complete!"
echo "========================"
echo ""

# Summary of key endpoints
echo "📋 Quick Reference:"
echo "🏠 Health Check: $BASE_URL/health"
echo "📖 API Root: $API_URL"
echo "🔐 Register: POST $API_URL/auth/register"
echo "🔑 Login: POST $API_URL/auth/login"
echo "👤 Profile: GET $API_URL/auth/profile"
echo "📝 Posts: GET $API_URL/posts/published"
echo "🔍 Search: GET $API_URL/posts/search?q=query"
echo "✍️  Create Post: POST $API_URL/posts"
echo ""
echo "💡 For detailed API documentation, see README.md"

# Check if .env file exists
if [ ! -f ".env" ]; then
    echo ""
    print_warning "No .env file found. Copy .env.example to .env and configure your settings."
fi

echo ""
echo "🎉 Happy coding!"