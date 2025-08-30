package seeder

import (
	"fmt"
	"log"
	"time"

	"github.com/kaungmyathan22/golang-multiuser-blog/internal/config"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/models"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/utils"
	"gorm.io/gorm"
)

type Seeder struct {
	db *gorm.DB
}

func NewSeeder() *Seeder {
	return &Seeder{
		db: config.GetDB(),
	}
}

// RunSeeder runs all seeders
func (s *Seeder) RunSeeder(force bool) error {
	log.Println("🌱 Starting database seeding...")

	// Check if data already exists
	if !force {
		var userCount, postCount, tagCount int64
		s.db.Model(&models.User{}).Count(&userCount)
		s.db.Model(&models.Post{}).Count(&postCount)
		s.db.Model(&models.Tag{}).Count(&tagCount)

		if userCount > 1 || postCount > 0 || tagCount > 5 { // More than default admin and default tags
			log.Println("ℹ️  Sample data already exists. Use --force to reseed.")
			return nil
		}
	}

	// Seed in order due to dependencies
	if err := s.seedUsers(); err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}

	if err := s.seedTags(); err != nil {
		return fmt.Errorf("failed to seed tags: %w", err)
	}

	if err := s.seedPosts(); err != nil {
		return fmt.Errorf("failed to seed posts: %w", err)
	}

	if err := s.seedComments(); err != nil {
		return fmt.Errorf("failed to seed comments: %w", err)
	}

	log.Println("✅ Database seeding completed successfully!")
	return nil
}

// seedUsers creates sample users
func (s *Seeder) seedUsers() error {
	log.Println("👥 Seeding users...")

	users := []models.User{
		{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Username:  "johndoe",
			Password:  "password123", // Will be hashed by BeforeCreate hook
			Bio:       "Passionate software developer and blogger. Love sharing knowledge about web development and emerging technologies.",
			Avatar:    "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150",
			IsActive:  true,
			IsAdmin:   false,
		},
		{
			FirstName: "Jane",
			LastName:  "Smith",
			Email:     "jane.smith@example.com",
			Username:  "janesmith",
			Password:  "password123",
			Bio:       "Tech enthusiast and writer. Specialized in UI/UX design and front-end development. Always learning something new.",
			Avatar:    "https://images.unsplash.com/photo-1494790108755-2616b612b272?w=150",
			IsActive:  true,
			IsAdmin:   false,
		},
		{
			FirstName: "Mike",
			LastName:  "Johnson",
			Email:     "mike.johnson@example.com",
			Username:  "mikej",
			Password:  "password123",
			Bio:       "Backend developer with 8+ years of experience. Expert in Go, Python, and distributed systems architecture.",
			Avatar:    "https://images.unsplash.com/photo-1560250097-0b93528c311a?w=150",
			IsActive:  true,
			IsAdmin:   false,
		},
		{
			FirstName: "Sarah",
			LastName:  "Wilson",
			Email:     "sarah.wilson@example.com",
			Username:  "sarahw",
			Password:  "password123",
			Bio:       "Full-stack developer and tech blogger. Passionate about clean code, testing, and agile methodologies.",
			Avatar:    "https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?w=150",
			IsActive:  true,
			IsAdmin:   false,
		},
		{
			FirstName: "Alex",
			LastName:  "Chen",
			Email:     "alex.chen@example.com",
			Username:  "alexchen",
			Password:  "password123",
			Bio:       "DevOps engineer and cloud architect. Love automating everything and building scalable infrastructure.",
			Avatar:    "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150",
			IsActive:  true,
			IsAdmin:   false,
		},
	}

	for _, user := range users {
		var existingUser models.User
		result := s.db.Where("email = ?", user.Email).First(&existingUser)
		if result.Error != nil {
			if err := s.db.Create(&user).Error; err != nil {
				return err
			}
			log.Printf("✅ Created user: %s (%s)", user.Username, user.Email)
		}
	}

	return nil
}

// seedTags creates additional sample tags
func (s *Seeder) seedTags() error {
	log.Println("🏷️  Seeding additional tags...")

	additionalTags := []models.Tag{
		{Name: "JavaScript", Slug: "javascript", Description: "JavaScript programming and frameworks", Color: "#F7DF1E"},
		{Name: "Go", Slug: "go", Description: "Go programming language", Color: "#00ADD8"},
		{Name: "Python", Slug: "python", Description: "Python programming and data science", Color: "#3776AB"},
		{Name: "Web Development", Slug: "web-development", Description: "Web development tips and tutorials", Color: "#61DAFB"},
		{Name: "DevOps", Slug: "devops", Description: "DevOps practices and tools", Color: "#FF6B35"},
		{Name: "Machine Learning", Slug: "machine-learning", Description: "ML and AI related content", Color: "#FF6F00"},
		{Name: "Mobile", Slug: "mobile", Description: "Mobile app development", Color: "#A4C639"},
		{Name: "Database", Slug: "database", Description: "Database design and optimization", Color: "#336791"},
		{Name: "Security", Slug: "security", Description: "Cybersecurity and best practices", Color: "#DC143C"},
		{Name: "Open Source", Slug: "open-source", Description: "Open source projects and contributions", Color: "#28A745"},
	}

	for _, tag := range additionalTags {
		var existingTag models.Tag
		result := s.db.Where("slug = ?", tag.Slug).First(&existingTag)
		if result.Error != nil {
			if err := s.db.Create(&tag).Error; err != nil {
				return err
			}
			log.Printf("✅ Created tag: %s", tag.Name)
		}
	}

	return nil
}

// seedPosts creates sample blog posts
func (s *Seeder) seedPosts() error {
	log.Println("📝 Seeding blog posts...")

	// Get users and tags for associations
	var users []models.User
	s.db.Where("email != ?", "admin@blog.com").Find(&users)

	var tags []models.Tag
	s.db.Find(&tags)

	if len(users) == 0 || len(tags) == 0 {
		return fmt.Errorf("users or tags not found for seeding posts")
	}

	posts := []struct {
		Title       string
		Content     string
		Status      models.PostStatus
		AuthorIndex int
		TagIndices  []int
	}{
		{
			Title: "Getting Started with Go: A Beginner's Guide",
			Content: `# Getting Started with Go

Go, also known as Golang, is a programming language developed by Google. It's designed for simplicity, efficiency, and reliability.

## Why Choose Go?

1. **Simple syntax** - Easy to learn and read
2. **Fast compilation** - Quick build times
3. **Concurrency support** - Built-in goroutines
4. **Strong standard library** - Comprehensive packages

## Hello World Example

` + "```go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Println(\"Hello, World!\")\n}\n```" + `

## Conclusion

Go is an excellent choice for modern software development, especially for web services and system programming.`,
			Status:      models.PostStatusPublished,
			AuthorIndex: 0,
			TagIndices:  []int{1, 2}, // Go, Tutorial
		},
		{
			Title: "Building REST APIs with Go and Gin",
			Content: `# Building REST APIs with Go and Gin

The Gin framework makes it easy to build high-performance REST APIs in Go.

## Setting Up Gin

First, install the Gin package:

` + "```bash\ngo get github.com/gin-gonic/gin\n```" + `

## Creating Your First API

` + "```go\npackage main\n\nimport \"github.com/gin-gonic/gin\"\n\nfunc main() {\n    r := gin.Default()\n    r.GET(\"/ping\", func(c *gin.Context) {\n        c.JSON(200, gin.H{\"message\": \"pong\"})\n    })\n    r.Run()\n}\n```" + `

## Best Practices

- Use middleware for common functionality
- Implement proper error handling
- Add request validation
- Use structured logging

This approach will help you build robust and maintainable APIs.`,
			Status:      models.PostStatusPublished,
			AuthorIndex: 2,
			TagIndices:  []int{1, 3, 2}, // Go, Web Development, Tutorial
		},
		{
			Title: "Modern JavaScript: ES6+ Features You Should Know",
			Content: `# Modern JavaScript: ES6+ Features

JavaScript has evolved significantly with ES6 and later versions. Here are the key features every developer should know.

## Arrow Functions

` + "```javascript\n// Traditional function\nfunction add(a, b) {\n    return a + b;\n}\n\n// Arrow function\nconst add = (a, b) => a + b;\n```" + `

## Destructuring

` + "```javascript\n// Array destructuring\nconst [first, second] = [1, 2];\n\n// Object destructuring\nconst {name, age} = person;\n```" + `

## Template Literals

` + "```javascript\nconst message = `Hello, ${name}! You are ${age} years old.`;\n```" + `

## Async/Await

` + "```javascript\nasync function fetchData() {\n    try {\n        const response = await fetch('/api/data');\n        const data = await response.json();\n        return data;\n    } catch (error) {\n        console.error('Error:', error);\n    }\n}\n```" + `

These features make JavaScript more powerful and easier to work with.`,
			Status:      models.PostStatusPublished,
			AuthorIndex: 1,
			TagIndices:  []int{0, 3, 2}, // JavaScript, Web Development, Tutorial
		},
		{
			Title: "DevOps Best Practices for Small Teams",
			Content: `# DevOps Best Practices for Small Teams

Implementing DevOps practices doesn't require a large team or complex infrastructure.

## Start with the Basics

1. **Version Control** - Use Git effectively
2. **Automated Testing** - Write and run tests automatically
3. **Continuous Integration** - Automate builds and deployments
4. **Monitoring** - Know when things go wrong

## Essential Tools

- **Git** for version control
- **Docker** for containerization
- **GitHub Actions** or **GitLab CI** for CI/CD
- **Prometheus** and **Grafana** for monitoring

## Implementation Strategy

Start small and gradually add more sophisticated practices:

1. Set up basic CI/CD pipeline
2. Add automated testing
3. Implement monitoring and alerting
4. Introduce infrastructure as code

## Common Pitfalls

- Trying to do too much at once
- Ignoring security considerations
- Not involving the entire team
- Focusing only on tools, not culture

Remember: DevOps is about culture and collaboration, not just tools.`,
			Status:      models.PostStatusPublished,
			AuthorIndex: 4,
			TagIndices:  []int{4, 0}, // DevOps, Technology
		},
		{
			Title: "Introduction to Machine Learning with Python",
			Content: `# Introduction to Machine Learning with Python

Machine Learning is transforming how we solve complex problems. Python makes it accessible to everyone.

## Why Python for ML?

- **Rich ecosystem** - NumPy, Pandas, Scikit-learn
- **Easy to learn** - Simple syntax
- **Great community** - Extensive documentation and tutorials
- **Powerful libraries** - TensorFlow, PyTorch for deep learning

## Getting Started

` + "```python\nimport pandas as pd\nfrom sklearn.model_selection import train_test_split\nfrom sklearn.linear_model import LinearRegression\n\n# Load data\ndata = pd.read_csv('data.csv')\n\n# Prepare features and target\nX = data[['feature1', 'feature2']]\ny = data['target']\n\n# Split data\nX_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2)\n\n# Train model\nmodel = LinearRegression()\nmodel.fit(X_train, y_train)\n\n# Make predictions\npredictions = model.predict(X_test)\n```" + `

## Key Concepts

1. **Supervised Learning** - Learning from labeled data
2. **Unsupervised Learning** - Finding patterns in unlabeled data
3. **Feature Engineering** - Creating meaningful input variables
4. **Model Evaluation** - Measuring performance

This is just the beginning of an exciting journey into ML!`,
			Status:      models.PostStatusPublished,
			AuthorIndex: 3,
			TagIndices:  []int{2, 5, 0}, // Python, Machine Learning, Technology
		},
		{
			Title: "Database Design Principles and Best Practices",
			Content: `# Database Design Principles and Best Practices

Good database design is crucial for application performance and maintainability.

## Normalization

Organize data to reduce redundancy:

1. **First Normal Form (1NF)** - Eliminate repeating groups
2. **Second Normal Form (2NF)** - Remove partial dependencies
3. **Third Normal Form (3NF)** - Remove transitive dependencies

## Indexing Strategy

` + "```sql\n-- Create index for frequently queried columns\nCREATE INDEX idx_user_email ON users(email);\n\n-- Composite index for multi-column queries\nCREATE INDEX idx_post_author_status ON posts(author_id, status);\n```" + `

## Performance Tips

- Use appropriate data types
- Avoid SELECT * queries
- Implement proper indexing
- Consider partitioning for large tables
- Use connection pooling

## Security Considerations

- Always use parameterized queries
- Implement proper authentication
- Apply principle of least privilege
- Regular security audits
- Encrypt sensitive data

## Common Mistakes

- Over-normalization
- Missing foreign key constraints
- Ignoring query performance
- Poor naming conventions
- Not planning for scalability

Proper database design pays dividends in the long run.`,
			Status:      models.PostStatusDraft,
			AuthorIndex: 2,
			TagIndices:  []int{7, 0}, // Database, Technology
		},
	}

	for i, postData := range posts {
		// Check if post already exists
		var existingPost models.Post
		slug := utils.GenerateSlug(postData.Title)
		result := s.db.Where("slug = ?", slug).First(&existingPost)

		if result.Error != nil {
			// Create post
			post := models.Post{
				Title:    postData.Title,
				Slug:     slug,
				Content:  postData.Content,
				Excerpt:  utils.ExtractExcerpt(postData.Content, 200),
				Status:   postData.Status,
				AuthorID: users[postData.AuthorIndex].ID,
			}

			// Set published date for published posts
			if postData.Status == models.PostStatusPublished {
				publishedTime := time.Now().AddDate(0, 0, -(len(posts) - i)) // Stagger dates
				post.PublishedAt = &publishedTime
			}

			if err := s.db.Create(&post).Error; err != nil {
				return err
			}

			// Add tags
			if len(postData.TagIndices) > 0 {
				var postTags []models.Tag
				for _, tagIndex := range postData.TagIndices {
					if tagIndex < len(tags) {
						postTags = append(postTags, tags[tagIndex])
					}
				}
				if err := s.db.Model(&post).Association("Tags").Append(&postTags); err != nil {
					log.Printf("Warning: Failed to add tags to post %s", post.Title)
				}
			}

			log.Printf("✅ Created post: %s (by %s)", post.Title, users[postData.AuthorIndex].Username)
		}
	}

	return nil
}

// seedComments creates sample comments
func (s *Seeder) seedComments() error {
	log.Println("💬 Seeding comments...")

	// Get published posts and users
	var posts []models.Post
	s.db.Where("status = ?", models.PostStatusPublished).Find(&posts)

	var users []models.User
	s.db.Where("email != ?", "admin@blog.com").Find(&users)

	if len(posts) == 0 || len(users) == 0 {
		return fmt.Errorf("posts or users not found for seeding comments")
	}

	comments := []struct {
		Content     string
		PostIndex   int
		AuthorIndex int
		Status      models.CommentStatus
	}{
		{
			Content:     "Great article! This really helped me understand the basics of Go. Looking forward to more content like this.",
			PostIndex:   0,
			AuthorIndex: 1,
			Status:      models.CommentStatusApproved,
		},
		{
			Content:     "Thanks for sharing this tutorial. The code examples are very clear and easy to follow.",
			PostIndex:   0,
			AuthorIndex: 2,
			Status:      models.CommentStatusApproved,
		},
		{
			Content:     "Excellent guide! I've been looking for a comprehensive introduction to Go. This covers all the essentials.",
			PostIndex:   0,
			AuthorIndex: 3,
			Status:      models.CommentStatusApproved,
		},
		{
			Content:     "The Gin framework looks very promising. Have you compared it with other Go web frameworks like Echo or Fiber?",
			PostIndex:   1,
			AuthorIndex: 0,
			Status:      models.CommentStatusApproved,
		},
		{
			Content:     "Nice tutorial! Could you also cover authentication and middleware in a follow-up post?",
			PostIndex:   1,
			AuthorIndex: 3,
			Status:      models.CommentStatusApproved,
		},
		{
			Content:     "ES6+ features have really improved JavaScript development. Arrow functions and destructuring are game-changers!",
			PostIndex:   2,
			AuthorIndex: 4,
			Status:      models.CommentStatusApproved,
		},
		{
			Content:     "Async/await has made working with promises so much cleaner. Great explanation!",
			PostIndex:   2,
			AuthorIndex: 0,
			Status:      models.CommentStatusApproved,
		},
		{
			Content:     "Very practical advice for small teams. We implemented CI/CD following similar principles and it made a huge difference.",
			PostIndex:   3,
			AuthorIndex: 1,
			Status:      models.CommentStatusApproved,
		},
		{
			Content:     "Docker has been a game-changer for our deployment process. Thanks for including it in the essential tools list.",
			PostIndex:   3,
			AuthorIndex: 2,
			Status:      models.CommentStatusPending,
		},
		{
			Content:     "Python's ML ecosystem is indeed amazing. Scikit-learn makes it so easy to get started with machine learning.",
			PostIndex:   4,
			AuthorIndex: 1,
			Status:      models.CommentStatusApproved,
		},
	}

	for _, commentData := range comments {
		if commentData.PostIndex < len(posts) && commentData.AuthorIndex < len(users) {
			comment := models.Comment{
				Content:  commentData.Content,
				Status:   commentData.Status,
				PostID:   posts[commentData.PostIndex].ID,
				AuthorID: users[commentData.AuthorIndex].ID,
			}

			if err := s.db.Create(&comment).Error; err != nil {
				return err
			}

			log.Printf("✅ Created comment on post: %s", posts[commentData.PostIndex].Title)
		}
	}

	// Create some reply comments
	var parentComments []models.Comment
	s.db.Where("parent_id IS NULL").Limit(3).Find(&parentComments)

	replies := []struct {
		Content     string
		ParentIndex int
		AuthorIndex int
	}{
		{
			Content:     "I agree! Go's simplicity is one of its strongest features. The learning curve is much gentler compared to other languages.",
			ParentIndex: 0,
			AuthorIndex: 4,
		},
		{
			Content:     "You should definitely check out Echo and Fiber. Each has its own strengths, but Gin has excellent documentation and community support.",
			ParentIndex: 1,
			AuthorIndex: 2,
		},
	}

	for _, replyData := range replies {
		if replyData.ParentIndex < len(parentComments) && replyData.AuthorIndex < len(users) {
			reply := models.Comment{
				Content:  replyData.Content,
				Status:   models.CommentStatusApproved,
				PostID:   parentComments[replyData.ParentIndex].PostID,
				AuthorID: users[replyData.AuthorIndex].ID,
				ParentID: &parentComments[replyData.ParentIndex].ID,
			}

			if err := s.db.Create(&reply).Error; err != nil {
				return err
			}

			log.Printf("✅ Created reply comment")
		}
	}

	return nil
}

// CleanDatabase removes all seeded data (useful for testing)
func (s *Seeder) CleanDatabase() error {
	log.Println("🧹 Cleaning database...")

	// Delete in reverse order of dependencies
	if err := s.db.Exec("DELETE FROM post_tags").Error; err != nil {
		return err
	}

	if err := s.db.Where("email != ?", "admin@blog.com").Delete(&models.Comment{}).Error; err != nil {
		return err
	}

	if err := s.db.Where("author_id != (SELECT id FROM users WHERE email = 'admin@blog.com' LIMIT 1)").Delete(&models.Post{}).Error; err != nil {
		return err
	}

	if err := s.db.Where("slug NOT IN ('technology', 'lifestyle', 'tutorial', 'news', 'opinion')").Delete(&models.Tag{}).Error; err != nil {
		return err
	}

	if err := s.db.Where("email != ?", "admin@blog.com").Delete(&models.User{}).Error; err != nil {
		return err
	}

	log.Println("✅ Database cleaned successfully")
	return nil
}
