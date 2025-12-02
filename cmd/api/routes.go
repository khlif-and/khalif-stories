package main

import (
	"time"

	"github.com/gin-gonic/gin"

	"khalif-stories/internal/config"
	"khalif-stories/pkg/middleware"

)

func SetupRoutes(r *gin.Engine, app *App, cfg *config.Config) {
	globalLimitConfig := middleware.RateLimitConfig{
		Limit:  300,
		Window: time.Minute,
	}
	r.Use(middleware.RateLimit(app.RDB, globalLimitConfig))

	authMiddleware := middleware.AuthMiddleware(cfg.JWTSecret)
	adminMiddleware := middleware.OnlyAdmin()

	r.GET("/api/categories", app.StoryHandler.GetCategories)
	r.GET("/api/categories/:id", app.StoryHandler.GetCategory)
	r.GET("/api/search/categories", app.SearchHandler.SearchCategories)
	r.GET("/api/stories", app.StoryHandler.GetAllStories)

	adminGroup := r.Group("/api/admin")
	adminGroup.Use(authMiddleware, adminMiddleware)
	{
		adminGroup.POST("/categories", app.StoryHandler.CreateCategory)
		adminGroup.PUT("/categories/:id", app.StoryHandler.UpdateCategory)
		adminGroup.DELETE("/categories/:id", app.StoryHandler.DeleteCategory)

		adminGroup.POST("/stories", app.StoryHandler.CreateStory)
		adminGroup.DELETE("/stories/:id", app.StoryHandler.DeleteStory)

		adminGroup.POST("/stories/:id/slides", app.ChapterHandler.AddSlide)
	}
}