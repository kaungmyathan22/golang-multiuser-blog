# Quick Start Guide

## 🚀 Get Started in 5 Minutes

### Prerequisites
- Go 1.21+
- PostgreSQL
- Air (optional, for hot reloading)

### 1. Setup Database
```sql
CREATE DATABASE golang_multiuser_blog;
```

### 2. Configure Environment
```bash
cp .env.example .env
# Edit .env with your database credentials
```

### 3. Run the Application
```bash
# Install dependencies
go mod tidy

# For development (with hot reloading)
air

# Or build and run
go build -o blog-server ./cmd/server
./blog-server
```

### 4. Test the API
```bash
# Make the test script executable and run it
chmod +x test_api.sh
./test_api.sh
```

## 🎯 Default Admin Account
- **Email**: admin@blog.com
- **Password**: admin123456
- **Username**: admin

⚠️ **Change this password immediately in production!**

## 📚 Quick API Examples

### Register a new user
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "username": "johndoe",
    "password": "password123"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email_or_username": "john@example.com",
    "password": "password123"
  }'
```

### Get published posts
```bash
curl http://localhost:8080/api/posts/published
```

### Create a post (requires authentication)
```bash
curl -X POST http://localhost:8080/api/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "My First Post",
    "content": "This is my first blog post!",
    "status": "published"
  }'
```

## 🔧 Key Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| POST | `/api/auth/register` | Register user |
| POST | `/api/auth/login` | Login user |
| GET | `/api/auth/profile` | Get user profile |
| GET | `/api/posts/published` | Get published posts |
| POST | `/api/posts` | Create post |
| GET | `/api/posts/{id}` | Get post by ID |
| GET | `/api/posts/search?q=term` | Search posts |

## 📖 Full Documentation
See [README.md](README.md) for complete documentation.

## 🆘 Troubleshooting

### Common Issues

**Database connection error?**
- Check PostgreSQL is running
- Verify database credentials in `.env`
- Ensure database exists

**Permission denied?**
- Check file permissions: `chmod +x blog-server`
- For test script: `chmod +x test_api.sh`

**Port already in use?**
- Change PORT in `.env`
- Or kill existing process: `lsof -ti:8080 | xargs kill`

**Build errors?**
- Run: `go mod tidy`
- Update Go to 1.21+

---

Happy coding! 🎉