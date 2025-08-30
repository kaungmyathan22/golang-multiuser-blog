#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${CYAN}🚀 Setting up Golang Multi-User Blog Development Environment${NC}"
echo ""

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check if Go is installed
if ! command_exists go; then
    echo -e "${RED}❌ Go is not installed. Please install Go first.${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Go is installed: $(go version)${NC}"

# Install Task (Taskfile runner)
echo -e "${YELLOW}📦 Installing Task runner...${NC}"
if command_exists task; then
    echo -e "${GREEN}✅ Task is already installed: $(task --version)${NC}"
else
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        if command_exists brew; then
            brew install go-task/tap/go-task
        else
            echo -e "${YELLOW}⚠️  Homebrew not found. Installing Task via Go...${NC}"
            go install github.com/go-task/task/v3/cmd/task@latest
        fi
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux
        go install github.com/go-task/task/v3/cmd/task@latest
    else
        # Other OS
        go install github.com/go-task/task/v3/cmd/task@latest
    fi

    if command_exists task; then
        echo -e "${GREEN}✅ Task installed successfully${NC}"
    else
        echo -e "${RED}❌ Failed to install Task${NC}"
    fi
fi

# Install Air (Hot reload tool)
echo -e "${YELLOW}🔥 Installing Air for hot reload...${NC}"
if command_exists air; then
    echo -e "${GREEN}✅ Air is already installed${NC}"
else
    go install github.com/cosmtrek/air@latest
    if command_exists air; then
        echo -e "${GREEN}✅ Air installed successfully${NC}"
    else
        echo -e "${RED}❌ Failed to install Air${NC}"
    fi
fi

# Install golangci-lint (Linter)
echo -e "${YELLOW}🔍 Installing golangci-lint...${NC}"
if command_exists golangci-lint; then
    echo -e "${GREEN}✅ golangci-lint is already installed${NC}"
else
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    if command_exists golangci-lint; then
        echo -e "${GREEN}✅ golangci-lint installed successfully${NC}"
    else
        echo -e "${RED}❌ Failed to install golangci-lint${NC}"
    fi
fi

# Create .gitignore if it doesn't exist
if [ ! -f .gitignore ]; then
    echo -e "${YELLOW}📝 Creating .gitignore file...${NC}"
    cat > .gitignore << EOF
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
# vendor/

# Go workspace file
go.work

# Build output
bin/
tmp/

# Air logs
*.log

# IDE files
.vscode/
.idea/
*.swp
*.swo
*~

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# Environment files
.env
.env.local
.env.development
.env.test
.env.production

# Database files
*.db
*.sqlite
*.sqlite3

# Coverage reports
coverage.out
coverage.html

# Docker
docker-compose.override.yml

# Kubernetes
*.kubeconfig
EOF
    echo -e "${GREEN}✅ .gitignore created${NC}"
else
    echo -e "${GREEN}✅ .gitignore already exists${NC}"
fi

echo ""
echo -e "${MAGENTA}🎉 Development environment setup complete!${NC}"
echo ""
echo -e "${CYAN}Available commands:${NC}"
echo -e "  ${GREEN}task dev${NC}          - Run server with hot reload"
echo -e "  ${GREEN}task build${NC}        - Build server"
echo -e "  ${GREEN}task test${NC}         - Run tests"
echo -e "  ${GREEN}task install${NC}      - Install dependencies"
echo ""
echo -e "${CYAN}Or use Make:${NC}"
echo -e "  ${GREEN}make help${NC}         - Show available make targets"
echo -e "  ${GREEN}make dev${NC}          - Run with hot reload"
echo -e "  ${GREEN}make build${NC}        - Build server"
echo -e "  ${GREEN}make install${NC}      - Install dependencies"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo -e "1. Run ${GREEN}task install${NC} or ${GREEN}make install${NC} to install Go dependencies"
echo -e "2. Start development with ${GREEN}task dev${NC} or ${GREEN}make dev${NC}"
echo -e "3. Check available commands with ${GREEN}task${NC} or ${GREEN}make help${NC}"
echo ""