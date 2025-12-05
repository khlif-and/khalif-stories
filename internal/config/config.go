package config

import (
	"log"

	"github.com/spf13/viper"

)

type Config struct {
	// Nama field disesuaikan dengan provider.go kamu
	DBUrl          string `mapstructure:"DATABASE_URL"`
	RedisAddr      string `mapstructure:"REDIS_ADDR"`
	Port           string `mapstructure:"PORT"`
	JWTSecret      string `mapstructure:"JWT_SECRET"`
	
	AzureConnStr   string `mapstructure:"AZURE_STORAGE_CONNECTION_STRING"`
	AzureContainer string `mapstructure:"AZURE_CONTAINER_NAME"`
	
	// Field Baru untuk Stories
	AzureContainerStoriesName string `mapstructure:"AZURE_CONTAINER_STORIES_NAME"`
	
	SlideLimit     int    `mapstructure:"SLIDE_LIMIT"`
}

// Kembalikan ke signature lama (1 return value) agar wire_gen.go tidak error
func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Jika error baca config, kita log saja (jangan return error biar signature tetap sama)
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Warning: gagal baca .env (pastikan file ada di root), menggunakan environment variables OS")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("Gagal parsing config:", err)
	}

	return &config
}