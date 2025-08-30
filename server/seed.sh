#!/bin/bash

# Seeder Build and Run Script

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    print_warning "This script should be run from the server directory"
    print_info "Current directory: $(pwd)"
    print_info "Looking for go.mod file..."
    exit 1
fi

# Parse command line arguments
COMMAND=""
FORCE=""
CLEAN=""
HELP=""

while [[ $# -gt 0 ]]; do
    case $1 in
        --seed|-s)
            COMMAND="seed"
            shift
            ;;
        --force|-f)
            FORCE="--force"
            shift
            ;;
        --clean|-c)
            COMMAND="clean"
            shift
            ;;
        --help|-h)
            HELP="true"
            shift
            ;;
        *)
            print_warning "Unknown option: $1"
            shift
            ;;
    esac
done

# Show help if requested
if [ "$HELP" = "true" ]; then
    echo "Seeder CLI Tool"
    echo "==============="
    echo ""
    echo "Usage:"
    echo "  ./seed.sh --seed          Seed the database with sample data"
    echo "  ./seed.sh --seed --force  Force reseeding even if data exists"
    echo "  ./seed.sh --clean         Clean the database (remove seeded data)"
    echo "  ./seed.sh --help          Show this help message"
    echo ""
    echo "Examples:"
    echo "  ./seed.sh --seed"
    echo "  ./seed.sh --seed --force"
    echo "  ./seed.sh --clean"
    exit 0
fi

# Build the seeder
print_info "Building seeder..."
go build -o seeder ./cmd/seeder
BUILD_RESULT=$?
if [ $BUILD_RESULT -ne 0 ]; then
    print_status 1 "Failed to build seeder"
    exit 1
fi
print_status 0 "Seeder built successfully"

# Run the seeder with the specified command
if [ "$COMMAND" = "seed" ]; then
    print_info "Seeding database..."
    ./seeder --seed $FORCE
    SEED_RESULT=$?
    if [ $SEED_RESULT -ne 0 ]; then
        print_status 1 "Failed to seed database"
        exit 1
    fi
    print_status 0 "Database seeded successfully"
elif [ "$COMMAND" = "clean" ]; then
    print_info "Cleaning database..."
    ./seeder --clean
    CLEAN_RESULT=$?
    if [ $CLEAN_RESULT -ne 0 ]; then
        print_status 1 "Failed to clean database"
        exit 1
    fi
    print_status 0 "Database cleaned successfully"
else
    print_info "Seeder built successfully. Run with --seed or --clean"
    echo ""
    echo "Usage:"
    echo "  ./seed.sh --seed          Seed the database with sample data"
    echo "  ./seed.sh --seed --force  Force reseeding even if data exists"
    echo "  ./seed.sh --clean         Clean the database (remove seeded data)"
    echo "  ./seed.sh --help          Show help message"
fi

# Clean up
rm -f seeder