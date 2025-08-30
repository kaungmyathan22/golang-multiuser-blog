package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/config"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/handlers"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/middleware"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/repository"
	"github.com/kaungmyathan22/golang-multiuser-blog/internal/service"
)

type Router struct {
	config         *config.Config
	authHandler    *handlers.AuthHandler
	postHandler    *handlers.PostHandler
	tagHandler     *handlers.TagHandler
	commentHandler *handlers.CommentHandler
	adminHandler   *handlers.AdminHandler
}

func NewRouter(cfg *config.Config) *Router {
	// Get database instance
	db := config.GetDB()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	tagRepo := repository.NewTagRepository(db)
	commentRepo := repository.NewCommentRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo, cfg)
	postService := service.NewPostService(postRepo, tagRepo, commentRepo)
	tagService := service.NewTagService(tagRepo)
	commentService := service.NewCommentService(commentRepo, postRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userService)
	postHandler := handlers.NewPostHandler(postService)
	tagHandler := handlers.NewTagHandler(tagService)
	commentHandler := handlers.NewCommentHandler(commentService)
	adminHandler := handlers.NewAdminHandler(userService)

	return &Router{
		config:         cfg,
		authHandler:    authHandler,
		postHandler:    postHandler,
		tagHandler:     tagHandler,
		commentHandler: commentHandler,
		adminHandler:   adminHandler,
	}
}

func (r *Router) SetupRoutes() *gin.Engine {
	// Set gin mode
	gin.SetMode(r.config.GinMode)

	// Create router
	router := gin.New()

	// Add middlewares
	router.Use(middleware.CORS())
	router.Use(middleware.RequestLoggerMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Server is running",
		})
	})

	// API routes
	api := router.Group("/api")
	{
		// Public routes (no authentication required)
		public := api.Group("")
		public.Use(middleware.PaginationMiddleware())
		{
			// Authentication routes
			auth := public.Group("/auth")
			{
				auth.POST("/register", r.authHandler.Register)
				auth.POST("/login", r.authHandler.Login)
				auth.POST("/refresh", r.authHandler.RefreshToken)
			}

			// Public post routes
			posts := public.Group("/posts")
			posts.Use(middleware.OptionalAuthMiddleware(r.config))
			{
				posts.GET("", r.postHandler.GetPosts)
				posts.GET("/published", r.postHandler.GetPublishedPosts)
				posts.GET("/search", r.postHandler.SearchPosts)
				posts.GET("/:id", r.postHandler.GetPost)
				posts.GET("/slug/:slug", r.postHandler.GetPostBySlug)
			}

			// Public tag routes
			tags := public.Group("/tags")
			tags.Use(middleware.OptionalAuthMiddleware(r.config))
			{
				tags.GET("", r.tagHandler.GetTags)
				tags.GET("/all", r.tagHandler.GetAllTags)
				tags.GET("/popular", r.tagHandler.GetPopularTags)
				tags.GET("/:id", r.tagHandler.GetTag)
				tags.GET("/slug/:slug", r.tagHandler.GetTagBySlug)
				tags.GET("/:id/posts", r.tagHandler.GetPostsByTag)
			}

			// Public comment routes (separate from posts to avoid conflicts)

			comments := public.Group("/comments")
			comments.Use(middleware.OptionalAuthMiddleware(r.config))
			{
				comments.GET("/post/:post_id", r.commentHandler.GetCommentsByPost)
			}
		}

		// Protected routes (authentication required)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(r.config))
		protected.Use(middleware.PaginationMiddleware())
		{
			// Protected auth routes
			auth := protected.Group("/auth")
			{
				auth.GET("/profile", r.authHandler.GetProfile)
				auth.PUT("/profile", r.authHandler.UpdateProfile)
				auth.POST("/change-password", r.authHandler.ChangePassword)
			}

			// Protected post routes
			posts := protected.Group("/posts")
			{
				posts.POST("", r.postHandler.CreatePost)
				posts.PUT("/:id", r.postHandler.UpdatePost)
				posts.DELETE("/:id", r.postHandler.DeletePost)
				posts.POST("/:id/publish", r.postHandler.PublishPost)
				posts.POST("/:id/unpublish", r.postHandler.UnpublishPost)
			}

			// Protected comment routes
			comments := protected.Group("/comments")
			{
				comments.POST("", r.commentHandler.CreateComment)
				comments.PUT("/:id", r.commentHandler.UpdateComment)
				comments.DELETE("/:id", r.commentHandler.DeleteComment)
				comments.GET("/my-comments", r.commentHandler.GetCommentsByAuthor)
			}
		}

		// Admin routes (admin access required)
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(r.config))
		admin.Use(middleware.AdminMiddleware())
		admin.Use(middleware.PaginationMiddleware())
		{
			// Admin user management
			adminUsers := admin.Group("/users")
			{
				adminUsers.GET("", r.adminHandler.GetUsers)
				adminUsers.GET("/:id", r.adminHandler.GetUser)
				adminUsers.POST("/:id/deactivate", r.adminHandler.DeactivateUser)
				adminUsers.POST("/:id/activate", r.adminHandler.ActivateUser)
				adminUsers.GET("/stats", r.adminHandler.GetUserStats)
			}

			// Admin post management
			adminPosts := admin.Group("/posts")
			{
				adminPosts.GET("", r.postHandler.GetPosts)
				adminPosts.GET("/:id", r.postHandler.GetPost)
				adminPosts.PUT("/:id", r.postHandler.UpdatePost)
				adminPosts.DELETE("/:id", r.postHandler.DeletePost)
				adminPosts.POST("/:id/publish", r.postHandler.PublishPost)
				adminPosts.POST("/:id/unpublish", r.postHandler.UnpublishPost)
			}

			// Admin comment management
			adminComments := admin.Group("/comments")
			{
				adminComments.GET("/pending", r.commentHandler.GetPendingComments)
				adminComments.POST("/:id/approve", r.commentHandler.ApproveComment)
				adminComments.POST("/:id/reject", r.commentHandler.RejectComment)
				adminComments.GET("/pending/count", r.commentHandler.GetPendingCount)
			}

			// Admin tag management
			adminTags := admin.Group("/tags")
			{
				adminTags.POST("", r.tagHandler.CreateTag)
				adminTags.PUT("/:id", r.tagHandler.UpdateTag)
				adminTags.DELETE("/:id", r.tagHandler.DeleteTag)
				adminTags.GET("/stats", r.tagHandler.GetTagStats)
			}

			// Admin dashboard
			admin.GET("/dashboard/stats", r.adminHandler.GetDashboardStats)
		}
	}

	return router
}
