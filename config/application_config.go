package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	SSUAnnouncementURL string
}

func LoadConfig() *AppConfig {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: .env file not found")
	}

	return &AppConfig{
		SSUAnnouncementURL: os.Getenv("SSU_ANNOUNCEMENT_URL"),
	}
}
