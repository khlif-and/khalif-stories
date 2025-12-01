package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"khalif-stories/internal/config"
	"khalif-stories/internal/domain"
	"khalif-stories/pkg/database"

)

func main() {
	refreshFlag := flag.Bool("refresh", false, "Reset Database")
	flag.Parse()

	cfg := config.LoadConfig()

	app, err := InitializeApp()
	if err != nil {
		log.Fatal(err)
	}

	if *refreshFlag {
		fmt.Println("Resetting Schema...")
		database.ResetSchema(app.DB)
	}

	app.DB.AutoMigrate(&domain.Category{}, &domain.Story{}, &domain.Slide{})

	database.SeedCategories(app.DB)

	r := gin.Default()

	SetupRoutes(r, app, cfg)

	r.Run(":" + cfg.Port)
}