package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Root endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Golang Multi-User Blog API", "status": "running", "version": "1.0.0"}`)
	})

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status": "healthy"}`)
	})

	// API routes placeholder
	http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Users API endpoint", "method": "%s"}`, r.Method)
	})

	http.HandleFunc("/api/posts", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Posts API endpoint", "method": "%s"}`, r.Method)
	})

	http.HandleFunc("/api/auth", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Authentication API endpoint", "method": "%s"}`, r.Method)
	})

	log.Printf("🚀 Golang Multi-User Blog Server starting on port %s", port)
	log.Printf("📁 API Documentation: http://localhost:%s/", port)
	log.Printf("💚 Health Check: http://localhost:%s/health", port)
	log.Printf("👥 Users API: http://localhost:%s/api/users", port)
	log.Printf("📝 Posts API: http://localhost:%s/api/posts", port)
	log.Printf("🔐 Auth API: http://localhost:%s/api/auth", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
