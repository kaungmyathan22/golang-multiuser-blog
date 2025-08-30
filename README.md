# Golang Multi-User Blog Platform

A clean, scalable multi-user blogging platform built with Go using a monolith architecture with no third-party dependencies.

## 🚀 Quick Start

### Prerequisites

- Go 1.19 or higher
- Git

### Setup

1. **Clone and setup the project:**
   ```bash
   git clone <your-repo-url>
   cd golang-multiuser-blog
   ./setup.sh  # Install development tools
   ```

2. **Install dependencies:**
   ```bash
   task install
   # or
   make install
   ```

3. **Start development server with hot reload:**
   ```bash
   task dev
   # or
   make dev
   ```

4. **Visit your application:**
   - API: http://localhost:8080
   - Health Check: http://localhost:8080/health
   - Users API: http://localhost:8080/api/users
   - Posts API: http://localhost:8080/api/posts
   - Auth API: http://localhost:8080/api/auth

## 🛠️ Development Tools

### Task Runner (Taskfile)

We use [Task](https://taskfile.dev/) as our primary task runner. Available commands:

```bash
task                # Show all available tasks
task dev            # Run server with hot reload
task build          # Build server
task test           # Run tests
task test:watch     # Run tests with file watching
task lint           # Run linter
task fmt            # Format code
```

### Make (Alternative)

If you prefer Make over Task:

```bash
make help           # Show all available commands
make dev            # Run with hot reload
make build          # Build server
make test           # Run tests
make lint           # Run linter
```

### Hot Reload

Hot reload is powered by [Air](https://github.com/cosmtrek/air). Configuration file:

- `.air.toml` - Main server configuration

## 📁 Project Structure

```
golang-multiuser-blog/
├── server/
│   ├── cmd/
│   │   └── server/               # Application entry point
│   ├── internal/             # Private application code
│   │   ├── handlers/             # HTTP handlers
│   │   ├── services/             # Business logic
│   │   ├── models/               # Data models
│   │   ├── repository/           # Data access layer
│   │   ├── middleware/           # HTTP middleware
│   │   └── utils/                # Utility functions
│   ├── pkg/                  # Public packages
│   └── migrations/           # Database migrations
├── web/                      # Static files (optional)
├── docs/                     # Documentation
├── Taskfile.yml              # Task runner configuration
├── Makefile                  # Make configuration
└── .env.example              # Environment variables template
```

## 🏗️ Architecture

### Monolith Architecture

This project follows a clean monolith architecture with:

- **Domain-Driven Design**: Organized by business domains (users, posts, auth)
- **Layered Architecture**: Clear separation between handlers, services, and repositories
- **Dependency Injection**: Loose coupling between components
- **Standard Library Focus**: Minimal external dependencies

### Key Principles

1. **Simplicity**: Clean, readable code with minimal complexity
2. **Testability**: Easy to unit test and integration test
3. **Maintainability**: Well-organized code structure
4. **Performance**: Efficient use of Go's standard library
5. **Scalability**: Designed to handle growth within monolith boundaries

## 🔧 Configuration

Copy `.env.example` to `.env.development` and modify as needed:

```bash
cp .env.example .env.development
```

Key configuration options:
- **Database**: PostgreSQL, MySQL, or SQLite
- **Cache**: Redis (optional)
- **Logging**: JSON or text format
- **Authentication**: JWT-based
- **File Uploads**: Configurable size and types

## 🧪 Testing

```bash
# Run all tests
task test
# or
make test

# Run tests with coverage
task test:coverage
# or
make test-coverage

# Watch tests (re-run on file changes)
task test:watch
# or
make test-watch
```

## 📦 Building

```bash
# Build server
task build
# or
make build
```

## 🐳 Docker

```bash
# Build Docker image
task docker:build
# or
make docker-build

# Run with Docker
task docker:run
# or
make docker-run
```

## 🚀 Deployment

### Local Development
```bash
task dev  # Start with hot reload
```

### Production
```bash
task build           # Build server
./bin/golang-multiuser-blog  # Run server
```

## 📚 API Documentation

API endpoints:

- **GET /**: API information
- **GET /health**: Health check
- **GET/POST /api/users**: User management
- **GET/POST /api/posts**: Blog post management
- **POST /api/auth**: Authentication

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `task test`
5. Run linter: `task lint`
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Troubleshooting

### Hot Reload Not Working
1. Ensure Air is installed: `go install github.com/cosmtrek/air@latest`
2. Check Air configuration file (`.air.toml`)
3. Verify file permissions

### Build Errors
1. Run `task mod:tidy` or `make mod-tidy`
2. Check Go version compatibility
3. Verify all dependencies are available

### Port Conflicts
1. Check if port 8080 is available
2. Change PORT environment variable
3. Update configuration as needed

## 📞 Support

For questions and support:
1. Check the documentation in `docs/`
2. Review existing issues
3. Create a new issue with detailed information