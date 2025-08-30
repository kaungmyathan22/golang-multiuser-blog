package main

import (
	"flag"
	"log"

	"github.com/kaungmyathan22/golang-multiuser-blog/internal/config"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/migration"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/seeder"
)

func main() {
	// Define command line flags
	seed := flag.Bool("seed", false, "Seed the database with sample data")
	force := flag.Bool("force", false, "Force reseeding even if data exists")
	clean := flag.Bool("clean", false, "Clean the database (remove seeded data)")
	help := flag.Bool("help", false, "Show help")

	flag.Parse()

	// Show help if requested
	if *help {
		flag.Usage()
		return
	}

	// Load configuration
	log.Println("🔧 Loading configuration...")
	cfg := config.LoadConfig()

	// Initialize database
	log.Println("🗄️  Initializing database...")
	config.InitDatabase(cfg)

	// Run database migrations
	log.Println("📊 Running database migrations...")
	if err := migration.RunMigrations(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Handle clean command
	if *clean {
		log.Println("🧹 Cleaning database...")
		s := seeder.NewSeeder()
		if err := s.CleanDatabase(); err != nil {
			log.Fatal("Failed to clean database:", err)
		}
		log.Println("✅ Database cleaned successfully!")
		return
	}

	// Handle seed command
	if *seed {
		log.Println("🌱 Seeding database...")
		s := seeder.NewSeeder()
		if err := s.RunSeeder(*force); err != nil {
			log.Fatal("Failed to seed database:", err)
		}
		log.Println("✅ Database seeding completed!")
		return
	}

	// If no command specified, show usage
	flag.Usage()
}
