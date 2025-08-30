package utils

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/models"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// GenerateSlug creates a URL-friendly slug from a string
func GenerateSlug(text string) string {
	// Convert to lowercase
	slug := strings.ToLower(text)

	// Replace spaces and special characters with hyphens
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Remove leading/trailing hyphens
	slug = strings.Trim(slug, "-")

	// Limit length
	if len(slug) > 100 {
		slug = slug[:100]
		slug = strings.Trim(slug, "-")
	}

	return slug
}

// ValidateStruct validates a struct using struct tags
func ValidateStruct(s interface{}) []models.ValidationError {
	var validationErrors []models.ValidationError

	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, models.ValidationError{
				Field:   err.Field(),
				Tag:     err.Tag(),
				Value:   fmt.Sprintf("%v", err.Value()),
				Message: getValidationMessage(err),
			})
		}
	}

	return validationErrors
}

// getValidationMessage returns a user-friendly validation message
func getValidationMessage(err validator.FieldError) string {
	field := err.Field()
	tag := err.Tag()
	param := err.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, param)
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", field, param)
	case "alphanum":
		return fmt.Sprintf("%s must contain only alphanumeric characters", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, param)
	case "hexcolor":
		return fmt.Sprintf("%s must be a valid hex color", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

// IsValidSlug checks if a string is a valid slug format
func IsValidSlug(slug string) bool {
	if slug == "" {
		return false
	}

	// Check if slug contains only lowercase letters, numbers, and hyphens
	reg := regexp.MustCompile(`^[a-z0-9-]+$`)
	if !reg.MatchString(slug) {
		return false
	}

	// Check if slug doesn't start or end with hyphen
	if strings.HasPrefix(slug, "-") || strings.HasSuffix(slug, "-") {
		return false
	}

	// Check if slug doesn't contain consecutive hyphens
	if strings.Contains(slug, "--") {
		return false
	}

	return true
}

// TruncateText truncates text to specified length and adds ellipsis
func TruncateText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}

	// Find the last space before the limit
	lastSpace := strings.LastIndex(text[:maxLength], " ")
	if lastSpace == -1 {
		lastSpace = maxLength
	}

	return text[:lastSpace] + "..."
}

// SanitizeText removes extra whitespace and normalizes text
func SanitizeText(text string) string {
	// Replace multiple spaces with single space
	reg := regexp.MustCompile(`\s+`)
	text = reg.ReplaceAllString(text, " ")

	// Trim leading and trailing whitespace
	text = strings.TrimSpace(text)

	return text
}

// IsAlphaNumericWithSpaces checks if string contains only alphanumeric characters and spaces
func IsAlphaNumericWithSpaces(text string) bool {
	for _, char := range text {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != ' ' {
			return false
		}
	}
	return true
}

// ExtractExcerpt extracts excerpt from content
func ExtractExcerpt(content string, maxLength int) string {
	// Remove HTML tags (basic)
	reg := regexp.MustCompile(`<[^>]*>`)
	plainText := reg.ReplaceAllString(content, "")

	// Sanitize and truncate
	plainText = SanitizeText(plainText)
	return TruncateText(plainText, maxLength)
}

// CalculatePagination calculates pagination values
func CalculatePagination(page, perPage int, total int64) models.PaginationMeta {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}
	if perPage > 100 {
		perPage = 100
	}

	totalPages := int((total + int64(perPage) - 1) / int64(perPage))

	return models.PaginationMeta{
		Page:       page,
		PerPage:    perPage,
		Total:      int(total),
		TotalPages: totalPages,
	}
}
