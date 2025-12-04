package config

import (
	"log"

	"github.com/spf13/viper"

)

type Config struct {
	DBUrl          string `mapstructure:"DATABASE_URL"`
	RedisAddr      string `mapstructure:"REDIS_ADDR"`
	Port           string `mapstructure:"PORT"`
	JWTSecret      string `mapstructure:"JWT_SECRET"`
	AzureConnStr   string `mapstructure:"AZURE_STORAGE_CONNECTION_STRING"`
	AzureContainer string `mapstructure:"AZURE_CONTAINER_NAME"`
	SlideLimit     int    `mapstructure:"SLIDE_LIMIT"`
}

func LoadConfig() *Config {
	viper.AddConfigPath(".")
	viper.AddConfigPath("../..")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	viper.SetDefault("SLIDE_LIMIT", 20)

	_ = viper.ReadInConfig()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}

	return &config
}