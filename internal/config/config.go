package config

import (
	"log"

	"github.com/spf13/viper"

)

type Config struct {
	DBUrl                     string `mapstructure:"DATABASE_URL"`
	RedisAddr                 string `mapstructure:"REDIS_ADDR"`
	Port                      string `mapstructure:"PORT"`
	JWTSecret                 string `mapstructure:"JWT_SECRET"`
	AzureConnStr              string `mapstructure:"AZURE_STORAGE_CONNECTION_STRING"`
	AzureContainer            string `mapstructure:"AZURE_CONTAINER_NAME"`
	AzureContainerStoriesName string `mapstructure:"AZURE_CONTAINER_STORIES_NAME"`
	SlideLimit                int    `mapstructure:"SLIDE_LIMIT"`
	
	StoriesThumbPath string `mapstructure:"STORIES_THUMB_PATH"`
	StoriesSlidePath string `mapstructure:"STORIES_SLIDE_PATH"`
}

func LoadConfig() *Config {
	viper.AutomaticEnv()
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	
	viper.AddConfigPath(".")
	viper.AddConfigPath("../../")
	viper.AddConfigPath("../")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("Failed to parse config:", err)
	}

    // Set Default Values jika di .env tidak ada
    if config.StoriesThumbPath == "" {
        config.StoriesThumbPath = "stories/thumbnails/"
    }
    if config.StoriesSlidePath == "" {
        config.StoriesSlidePath = "stories/slides/"
    }

	if config.DBUrl == "" {
		log.Fatal("FATAL: DATABASE_URL is empty.")
	}

	return &config
}