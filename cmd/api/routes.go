package main

import (
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "khalif-stories/docs"
	"khalif-stories/internal/config"
	"khalif-stories/pkg/middleware"

)

func SetupRoutes(r *gin.Engine, app *App, cfg *config.Config) {
	r.Use(middleware.Logger())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	limiter := middleware.RateLimitConfig{Limit: 300, Window: time.Minute}
	r.Use(middleware.RateLimit(app.RDB, limiter))
	
	auth := middleware.AuthMiddleware(cfg.JWTSecret)
	admin := middleware.OnlyAdmin()

	r.GET("/api/categories", app.CategoryHandler.GetAll)
	r.GET("/api/categories/:id", app.CategoryHandler.GetOne)
	r.GET("/api/search/categories", app.CategoryHandler.Search)
	r.GET("/api/stories", app.StoryHandler.GetAll)
	r.GET("/api/stories/:uuid", app.StoryHandler.GetOne)
	r.GET("/api/search/stories", app.StoryHandler.Search)
	r.GET("/api/chapters/:uuid", app.ChapterHandler.GetOne)

	protected := r.Group("/api")
	protected.Use(auth)
	{
		protected.GET("/stories/recommendations", app.StoryHandler.GetRecommendations)
		protected.POST("/preferences", app.PreferenceHandler.Save)
	}

	adm := r.Group("/api/admin")
	adm.Use(auth, admin)
	{
		adm.POST("/categories", app.CategoryHandler.Create)
		adm.PUT("/categories/:id", app.CategoryHandler.Update)
		adm.DELETE("/categories/:id", app.CategoryHandler.Delete)
		adm.POST("/stories", app.StoryHandler.Create)
		adm.PUT("/stories/:uuid", app.StoryHandler.Update)
		adm.DELETE("/stories/:uuid", app.StoryHandler.Delete)
		adm.POST("/stories/:uuid/slides", app.StoryHandler.AddSlide)
		adm.POST("/chapters", app.ChapterHandler.Create)
		adm.DELETE("/chapters/:uuid", app.ChapterHandler.Delete)
		adm.POST("/chapters/:uuid/slides", app.ChapterHandler.AddSlide)
	}
}