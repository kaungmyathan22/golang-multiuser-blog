# Golang Multi-User Blog API

A comprehensive, production-ready multi-user blog API built with Go, PostgreSQL, and JWT authentication. This API provides complete functionality for user management, blog posts, comments, and tags with proper validation and security measures.

## 🚀 Features

### 🔐 Authentication & Authorization
- User registration and login
- JWT-based authentication
- Password hashing with bcrypt
- Role-based access control (Admin/User)
- Token refresh functionality
- Profile management

### 📝 Blog Management
- Create, read, update, delete blog posts
- Draft, published, and archived post statuses
- Automatic slug generation from titles
- Featured images and excerpts
- View count tracking
- Search functionality
- Tag-based categorization

### 💬 Comment System
- Nested comments (replies)
- Comment moderation (pending, approved, rejected)
- Author-based comment management

### 🏷️ Tag System
- Create and manage tags
- Color-coded tags
- Tag-based post filtering
- Popular tags functionality

### 🔧 Technical Features
- Clean architecture (Repository -> Service -> Handler)
- Comprehensive input validation
- Pagination support
- CORS middleware
- Request logging
- Error handling
- Database migrations
- Environment-based configuration

## 📋 Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Air (for hot reloading during development)

## 🛠️ Installation & Setup

### 1. Clone the Repository
```bash
git clone <repository-url>
cd golang-multiuser-blog/server
```

### 2. Install Dependencies
```bash
go mod download
```

### 3. Install Air for Hot Reloading
```bash
go install github.com/cosmtrek/air@latest
```

### 4. Setup PostgreSQL Database
```sql
-- Connect to PostgreSQL and create database
CREATE DATABASE golang_multiuser_blog;
CREATE USER blog_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE golang_multiuser_blog TO blog_user;
```

### 5. Environment Configuration
Copy the example environment file and configure your settings:
```bash
cp .env.example .env
```

Edit `.env` with your configuration:
```env
# Server Configuration
PORT=8080
GIN_MODE=debug

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=blog_user
DB_PASSWORD=your_password
DB_NAME=golang_multiuser_blog
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRES_IN=24h

# Application Configuration
APP_ENV=development
LOG_LEVEL=info
```

### 6. Run the Application

#### Development (with hot reloading)
```bash
air
```

#### Production
```bash
go build -o blog-server ./cmd/server
./blog-server
```

## 📚 API Documentation

### Base URL
```
http://localhost:8080/api
```

### Authentication Endpoints

#### Register User
```http
POST /api/auth/register
Content-Type: application/json

{
  \"first_name\": \"John\",
  \"last_name\": \"Doe\",
  \"email\": \"john@example.com\",
  \"username\": \"johndoe\",
  \"password\": \"password123\",
  \"bio\": \"Software developer and blogger\",
  \"avatar\": \"https://example.com/avatar.jpg\"
}
```

#### Login
```http
POST /api/auth/login
Content-Type: application/json

{
  \"email_or_username\": \"john@example.com\",
  \"password\": \"password123\"
}
```

#### Get Profile
```http
GET /api/auth/profile
Authorization: Bearer <your-jwt-token>
```

#### Update Profile
```http
PUT /api/auth/profile
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  \"first_name\": \"John\",
  \"last_name\": \"Smith\",
  \"bio\": \"Updated bio\"
}
```

#### Change Password
```http
POST /api/auth/change-password
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  \"old_password\": \"oldpassword\",
  \"new_password\": \"newpassword123\"
}
```

### Post Endpoints

#### Get Published Posts
```http
GET /api/posts/published?page=1&per_page=10
```

#### Get Post by ID
```http
GET /api/posts/{id}
```

#### Get Post by Slug
```http
GET /api/posts/slug/{slug}
```

#### Search Posts
```http
GET /api/posts/search?q=search-term&page=1&per_page=10
```

#### Create Post (Authenticated)
```http
POST /api/posts
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  \"title\": \"My First Blog Post\",
  \"content\": \"This is the content of my blog post...\",
  \"excerpt\": \"Short description\",
  \"featured_image\": \"https://example.com/image.jpg\",
  \"status\": \"published\",
  \"tag_ids\": [1, 2, 3]
}
```

#### Update Post (Authenticated)
```http
PUT /api/posts/{id}
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  \"title\": \"Updated Title\",
  \"content\": \"Updated content...\",
  \"status\": \"published\"
}
```

#### Delete Post (Authenticated)
```http
DELETE /api/posts/{id}
Authorization: Bearer <your-jwt-token>
```

#### Publish Post (Authenticated)
```http
POST /api/posts/{id}/publish
Authorization: Bearer <your-jwt-token>
```

#### Unpublish Post (Authenticated)
```http
POST /api/posts/{id}/unpublish
Authorization: Bearer <your-jwt-token>
```

## 🗂️ Project Structure

```
server/
├── cmd/server/                 # Application entry point
│   └── main.go
├── internal/                   # Private application code
│   ├── config/                # Configuration management
│   │   └── config.go
│   ├── handlers/              # HTTP handlers
│   │   ├── auth_handler.go
│   │   └── post_handler.go
│   ├── middleware/            # HTTP middleware
│   │   └── middleware.go
│   ├── migration/             # Database migrations
│   │   └── migration.go
│   ├── models/                # Data models
│   │   ├── user.go
│   │   ├── post.go
│   │   ├── comment.go
│   │   ├── tag.go
│   │   └── response.go
│   ├── repository/            # Data access layer
│   │   ├── user_repository.go
│   │   ├── post_repository.go
│   │   ├── comment_repository.go
│   │   └── tag_repository.go
│   ├── router/                # Route configuration
│   │   └── router.go
│   ├── service/               # Business logic layer
│   │   ├── user_service.go
│   │   ├── post_service.go
│   │   ├── comment_service.go
│   │   └── tag_service.go
│   └── utils/                 # Utility functions
│       ├── utils.go
│       └── jwt.go
├── .air.toml                  # Air configuration
├── .env.example              # Environment variables template
├── go.mod                    # Go modules
└── README.md                 # This file
```

## 🏗️ Architecture

This application follows Clean Architecture principles:

1. **Models Layer**: Defines data structures and business entities
2. **Repository Layer**: Handles data persistence and database operations
3. **Service Layer**: Contains business logic and validation
4. **Handler Layer**: Manages HTTP requests and responses
5. **Router Layer**: Configures routes and middleware

## 🔒 Security Features

- Password hashing using bcrypt
- JWT token-based authentication
- Input validation and sanitization
- SQL injection prevention with GORM
- CORS protection
- Rate limiting ready (can be added)
- Environment-based configuration

## 🗄️ Database Schema

### Users Table
- ID, FirstName, LastName, Email, Username
- Password (hashed), Bio, Avatar
- IsActive, IsAdmin, CreatedAt, UpdatedAt

### Posts Table
- ID, Title, Slug, Content, Excerpt
- FeaturedImg, Status, ViewCount, AuthorID
- PublishedAt, CreatedAt, UpdatedAt

### Comments Table
- ID, Content, Status, AuthorID, PostID
- ParentID (for nested comments)
- CreatedAt, UpdatedAt

### Tags Table
- ID, Name, Slug, Description, Color
- CreatedAt, UpdatedAt

### Post_Tags Table (Many-to-Many)
- PostID, TagID

## 🚀 Deployment

### Docker Deployment (Coming Soon)
```dockerfile
# Example Dockerfile structure
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o blog-server ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/blog-server .
CMD [\"./blog-server\"]
```

### Environment Variables for Production
```env
PORT=8080
GIN_MODE=release
DB_HOST=your-postgres-host
DB_USER=your-db-user
DB_PASSWORD=your-secure-password
JWT_SECRET=your-very-secure-jwt-secret
APP_ENV=production
```

## 🧪 Testing

### Run Tests
```bash
go test ./...
```

### Test with Air
```bash
air -c .air.test.toml
```

## 📈 Performance Considerations

- Database connection pooling configured
- Pagination implemented for all list endpoints
- Efficient database queries with proper indexing
- JWT token expiration and refresh mechanism
- View count increment optimized with goroutines

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License.

## 📞 Support

For support, email your-email@example.com or create an issue in the repository.

## 🚧 Roadmap

- [ ] Email verification
- [ ] Password reset functionality
- [ ] File upload for images
- [ ] Rate limiting
- [ ] API documentation with Swagger
- [ ] Comment management API
- [ ] Tag management API
- [ ] User management API (Admin)
- [ ] Analytics and statistics
- [ ] Search optimization
- [ ] Caching layer (Redis)
- [ ] Docker containerization
- [ ] CI/CD pipeline

---

## Default Admin Account

After running the application for the first time, a default admin account is created:

- **Email**: admin@blog.com
- **Password**: admin123456
- **Username**: admin

**⚠️ Important**: Change the default admin password immediately in production!

---

Built with ❤️ using Go, PostgreSQL, and modern web development practices.