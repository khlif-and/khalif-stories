package config

import (
	"log"
	"os"

	"github.com/spf13/viper"

)

type Config struct {
	DBUrl                       string `mapstructure:"DATABASE_URL"`
	RedisAddr                   string `mapstructure:"REDIS_ADDR"`
	Port                        string `mapstructure:"PORT"`
	JWTSecret                   string `mapstructure:"JWT_SECRET"`
	AzureConnStr                string `mapstructure:"AZURE_STORAGE_CONNECTION_STRING"`
	AzureContainer              string `mapstructure:"AZURE_CONTAINER_NAME"`
	AzureContainerStoriesName   string `mapstructure:"AZURE_CONTAINER_STORIES_NAME"`
	AzureContainerChapterImages string `mapstructure:"AZURE_CONTAINER_CHAPTER_IMAGES"`
	AzureContainerChapterSounds string `mapstructure:"AZURE_CONTAINER_CHAPTER_SOUNDS"`
	SlideLimit                  int    `mapstructure:"SLIDE_LIMIT"`
	StoriesThumbPath            string `mapstructure:"STORIES_THUMB_PATH"`
	StoriesSlidePath            string `mapstructure:"STORIES_SLIDE_PATH"`
}

func LoadConfig() *Config {
	viper.AutomaticEnv()
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AddConfigPath(".")
	viper.AddConfigPath("../../")
	viper.AddConfigPath("../")
	viper.AddConfigPath("/app")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Info: .env file not found, relying on System Environment Variables")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("Failed to parse config:", err)
	}

	if config.DBUrl == "" {
		config.DBUrl = os.Getenv("DATABASE_URL")
	}
	if config.RedisAddr == "" {
		config.RedisAddr = os.Getenv("REDIS_ADDR")
	}
	if config.Port == "" {
		config.Port = os.Getenv("PORT")
	}
	if config.JWTSecret == "" {
		config.JWTSecret = os.Getenv("JWT_SECRET")
	}
	if config.AzureConnStr == "" {
		config.AzureConnStr = os.Getenv("AZURE_STORAGE_CONNECTION_STRING")
	}
	if config.AzureContainer == "" {
		config.AzureContainer = os.Getenv("AZURE_CONTAINER_NAME")
	}
	if config.AzureContainerStoriesName == "" {
		config.AzureContainerStoriesName = os.Getenv("AZURE_CONTAINER_STORIES_NAME")
	}
	if config.AzureContainerChapterImages == "" {
		config.AzureContainerChapterImages = os.Getenv("AZURE_CONTAINER_CHAPTER_IMAGES")
	}
	if config.AzureContainerChapterSounds == "" {
		config.AzureContainerChapterSounds = os.Getenv("AZURE_CONTAINER_CHAPTER_SOUNDS")
	}
	if config.StoriesThumbPath == "" {
		config.StoriesThumbPath = "stories/thumbnails/"
	}
	if config.StoriesSlidePath == "" {
		config.StoriesSlidePath = "stories/slides/"
	}

	if config.DBUrl == "" {
		log.Fatal("FATAL: DATABASE_URL is empty. Please check your docker-compose.yml")
	}

	return &config
}