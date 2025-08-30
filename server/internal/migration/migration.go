package migration

import (
	"log"

	"github.com/kaungmyathan22/golang-multiuser-blog/internal/config"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/models"
)

// RunMigrations runs all database migrations
func RunMigrations() error {
	db := config.GetDB()

	log.Println("🔄 Running database migrations...")

	// Auto-migrate all models
	err := db.AutoMigrate(
		&models.User{},
		&models.Tag{},
		&models.Post{},
		&models.Comment{},
	)

	if err != nil {
		log.Printf("❌ Migration failed: %v", err)
		return err
	}

	log.Println("✅ Database migrations completed successfully")

	// Create default admin user if it doesn't exist
	if err := createDefaultAdmin(); err != nil {
		log.Printf("⚠️  Warning: Failed to create default admin user: %v", err)
	}

	// Create default tags if they don't exist
	if err := createDefaultTags(); err != nil {
		log.Printf("⚠️  Warning: Failed to create default tags: %v", err)
	}

	return nil
}

// createDefaultAdmin creates a default admin user
func createDefaultAdmin() error {
	db := config.GetDB()

	var user models.User
	result := db.Where("email = ?", "admin@blog.com").First(&user)

	if result.Error == nil {
		log.Println("ℹ️  Default admin user already exists")
		return nil
	}

	adminUser := models.User{
		FirstName: "Admin",
		LastName:  "User",
		Email:     "admin@blog.com",
		Username:  "admin",
		Password:  "admin123456", // This will be hashed by the BeforeCreate hook
		Bio:       "Default administrator account",
		IsActive:  true,
		IsAdmin:   true,
	}

	if err := db.Create(&adminUser).Error; err != nil {
		return err
	}

	log.Println("✅ Default admin user created (admin@blog.com / admin123456)")
	return nil
}

// createDefaultTags creates default tags
func createDefaultTags() error {
	db := config.GetDB()

	defaultTags := []models.Tag{
		{Name: "Technology", Slug: "technology", Description: "Posts about technology and programming", Color: "#3B82F6"},
		{Name: "Lifestyle", Slug: "lifestyle", Description: "Posts about lifestyle and personal experiences", Color: "#10B981"},
		{Name: "Tutorial", Slug: "tutorial", Description: "Educational and how-to posts", Color: "#F59E0B"},
		{Name: "News", Slug: "news", Description: "Latest news and updates", Color: "#EF4444"},
		{Name: "Opinion", Slug: "opinion", Description: "Personal opinions and thoughts", Color: "#8B5CF6"},
	}

	for _, tag := range defaultTags {
		var existingTag models.Tag
		result := db.Where("slug = ?", tag.Slug).First(&existingTag)

		if result.Error != nil {
			if err := db.Create(&tag).Error; err != nil {
				log.Printf("Failed to create tag %s: %v", tag.Name, err)
			} else {
				log.Printf("✅ Created default tag: %s", tag.Name)
			}
		}
	}

	return nil
}
