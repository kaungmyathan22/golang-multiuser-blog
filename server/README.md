# Golang Multi-User Blog Server

This is the backend server for the Golang Multi-User Blog application.

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose (for containerized deployment)
- PostgreSQL (if running without Docker)

## Getting Started

### Running with Docker (Recommended)

```bash
# Start the application with PostgreSQL database
docker-compose up --build

# Run in detached mode
docker-compose up --build -d

# Stop the application
docker-compose down
```

### Running locally

1. Set up environment variables:
   ```bash
   cp .env.example .env
   # Edit .env file with your configuration
   ```

2. Install dependencies:
   ```bash
   make deps
   ```

3. Run the application:
   ```bash
   make run
   ```

## Testing

### Unit Tests
```bash
make test-unit
```

### Integration Tests
```bash
make test-integration
```

### End-to-End Tests
```bash
make test-e2e
```

### All Tests
```bash
make test-all
```

### Test with Docker
```bash
./scripts/test.sh
```

## Building

### Build Binary
```bash
make build
```

### Build Docker Image
```bash
make docker-build
```

## Development

### Format Code
```bash
make fmt
```

### Vet Code
```bash
make vet
```

### Clean Build Files
```bash
make clean
```

## API Endpoints

- Health Check: `GET /health`
- Auth Endpoints:
  - Register: `POST /api/auth/register`
  - Login: `POST /api/auth/login`
  - Refresh Token: `POST /api/auth/refresh`
  - Get Profile: `GET /api/auth/profile`
  - Update Profile: `PUT /api/auth/profile`
  - Change Password: `POST /api/auth/change-password`

- Post Endpoints:
  - Get Posts: `GET /api/posts`
  - Get Published Posts: `GET /api/posts/published`
  - Search Posts: `GET /api/posts/search`
  - Get Post by ID: `GET /api/posts/:id`
  - Get Post by Slug: `GET /api/posts/slug/:slug`
  - Create Post: `POST /api/posts` (authenticated)
  - Update Post: `PUT /api/posts/:id` (authenticated)
  - Delete Post: `DELETE /api/posts/:id` (authenticated)
  - Publish Post: `POST /api/posts/:id/publish` (authenticated)
  - Unpublish Post: `POST /api/posts/:id/unpublish` (authenticated)

- Tag Endpoints:
  - Get Tags: `GET /api/tags`
  - Get All Tags: `GET /api/tags/all`
  - Get Popular Tags: `GET /api/tags/popular`
  - Get Tag by ID: `GET /api/tags/:id`
  - Get Tag by Slug: `GET /api/tags/slug/:slug`
  - Get Posts by Tag: `GET /api/tags/:id/posts`

- Comment Endpoints:
  - Get Comments by Post: `GET /api/comments/post/:post_id`
  - Create Comment: `POST /api/comments` (authenticated)
  - Update Comment: `PUT /api/comments/:id` (authenticated)
  - Delete Comment: `DELETE /api/comments/:id` (authenticated)
  - Get My Comments: `GET /api/comments/my-comments` (authenticated)

- Admin Endpoints:
  - Get Users: `GET /api/admin/users` (admin only)
  - Get User: `GET /api/admin/users/:id` (admin only)
  - Deactivate User: `POST /api/admin/users/:id/deactivate` (admin only)
  - Activate User: `POST /api/admin/users/:id/activate` (admin only)
  - Get User Stats: `GET /api/admin/users/stats` (admin only)
  - Get Pending Comments: `GET /api/admin/comments/pending` (admin only)
  - Approve Comment: `POST /api/admin/comments/:id/approve` (admin only)
  - Reject Comment: `POST /api/admin/comments/:id/reject` (admin only)
  - Get Pending Count: `GET /api/admin/comments/pending/count` (admin only)
  - Create Tag: `POST /api/admin/tags` (admin only)
  - Update Tag: `PUT /api/admin/tags/:id` (admin only)
  - Delete Tag: `DELETE /api/admin/tags/:id` (admin only)
  - Get Tag Stats: `GET /api/admin/tags/stats` (admin only)
  - Get Dashboard Stats: `GET /api/admin/dashboard/stats` (admin only)

## Environment Variables

See `.env.example` for all available environment variables.

## CI/CD

This project uses GitHub Actions for CI/CD:

- Unit tests
- Integration tests
- Build process
- Deployment (to be configured)

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
