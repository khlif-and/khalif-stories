package main

import (
	"time"

	"github.com/gin-gonic/gin"

	"khalif-stories/internal/config"
	"khalif-stories/pkg/middleware"

)

func SetupRoutes(r *gin.Engine, app *App, cfg *config.Config) {
	globalLimitConfig := middleware.RateLimitConfig{
		Limit:  60,
		Window: time.Minute,
	}
	r.Use(middleware.RateLimit(app.RDB, globalLimitConfig))

	authMiddleware := middleware.AuthMiddleware(cfg.JWTSecret)
	adminMiddleware := middleware.OnlyAdmin()

	adminGroup := r.Group("/api/admin")
	adminGroup.Use(authMiddleware, adminMiddleware)
	{
		adminGroup.POST("/categories", app.StoryHandler.CreateCategory)
		adminGroup.GET("/categories", app.StoryHandler.GetCategories)

		adminGroup.GET("/stories", app.StoryHandler.GetAllStories)
		adminGroup.POST("/stories", app.StoryHandler.CreateStory)
		adminGroup.DELETE("/stories/:id", app.StoryHandler.DeleteStory)

		adminGroup.POST("/stories/:id/slides", app.ChapterHandler.AddSlide)
	}
}