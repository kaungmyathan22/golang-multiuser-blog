package main

import (
	"log"
	"net/http"
	"time"

	"github.com/kaungmyathan22/golang-multiuser-blog/internal/config"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/migration"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/router"
)

func main() {
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

	// Initialize router
	log.Println("🌐 Setting up routes...")
	r := router.NewRouter(cfg)
	appRouter := r.SetupRoutes()

	// Configure server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      appRouter,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Print startup information
	log.Printf("🚀 Golang Multi-User Blog Server starting on port %s", cfg.Port)
	log.Printf("📱 Environment: %s", cfg.App.Environment)
	log.Printf("📄 Health Check: http://localhost:%s/health", cfg.Port)
	log.Printf("🔗 API Base URL: http://localhost:%s/api", cfg.Port)
	log.Printf("👤 Auth Endpoints:")
	log.Printf("   📝 Register: POST http://localhost:%s/api/auth/register", cfg.Port)
	log.Printf("   🔑 Login: POST http://localhost:%s/api/auth/login", cfg.Port)
	log.Printf("   👥 Profile: GET http://localhost:%s/api/auth/profile", cfg.Port)
	log.Printf("📝 Blog Endpoints:")
	log.Printf("   📚 Posts: GET http://localhost:%s/api/posts", cfg.Port)
	log.Printf("   📄 Published: GET http://localhost:%s/api/posts/published", cfg.Port)
	log.Printf("   🔍 Search: GET http://localhost:%s/api/posts/search?q=query", cfg.Port)
	log.Printf("💾 Database: PostgreSQL on %s:%d", cfg.Database.Host, cfg.Database.Port)
	log.Println("")
	log.Println("🎉 Server is ready to accept connections!")

	// Start server
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("❌ Server failed to start:", err)
	}
}
