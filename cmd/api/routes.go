package main

import (
	"time"

	"github.com/gin-gonic/gin"

	"khalif-stories/internal/config"
	"khalif-stories/pkg/middleware"

)

func SetupRoutes(r *gin.Engine, app *App, cfg *config.Config) {
	r.Use(middleware.Logger())

	limiter := middleware.RateLimitConfig{Limit: 300, Window: time.Minute}
	r.Use(middleware.RateLimit(app.RDB, limiter))

	auth := middleware.AuthMiddleware(cfg.JWTSecret)
	admin := middleware.OnlyAdmin()

	// Categories tetap pakai :id (karena biasanya pakai integer ID)
	r.GET("/api/categories", app.CategoryHandler.GetAll)
	r.GET("/api/categories/:id", app.CategoryHandler.GetOne)
	r.GET("/api/search/categories", app.CategoryHandler.Search)

	r.GET("/api/stories", app.StoryHandler.GetAll)
	r.GET("/api/search/stories", app.StoryHandler.Search)

	adm := r.Group("/api/admin")
	adm.Use(auth, admin)
	{
		adm.POST("/categories", app.CategoryHandler.Create)
		adm.PUT("/categories/:id", app.CategoryHandler.Update)
		adm.DELETE("/categories/:id", app.CategoryHandler.Delete)

		adm.POST("/stories", app.StoryHandler.Create)
		
		// PERBAIKAN DISINI: Ubah :id menjadi :uuid
		adm.PUT("/stories/:uuid", app.StoryHandler.Update)
		adm.DELETE("/stories/:uuid", app.StoryHandler.Delete)
		adm.POST("/stories/:uuid/slides", app.StoryHandler.AddSlide)
	}
}