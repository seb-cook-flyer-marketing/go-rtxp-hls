package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Configuration struct {
	Port   string
	Secret string
	URL    string
	FFmpeg string
}

var Config Configuration

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	goEnv := os.Getenv("ENV")
	if goEnv == "" {
		goEnv = "development"
	}

	switch goEnv {
	case "production":
		Config = getProductionConfig()
	case "staging":
		Config = getStagingConfig()
	default:
		Config = getDevelopmentConfig()
	}
}

func getDevelopmentConfig() Configuration {
	// Implement development configuration
	return Configuration{
		Port:   "8080",
		Secret: "dev-secret",
		URL:    "http://localhost:8080",
		FFmpeg: "/usr/local/bin/ffmpeg",
	}
}

func getStagingConfig() Configuration {
	// Implement staging configuration
	return Configuration{
		Port:   "8080",
		Secret: "dev-secret",
		URL:    "http://localhost:8080",
		FFmpeg: "/usr/local/bin/ffmpeg",
	}
}

func getProductionConfig() Configuration {
	// Implement production configuration
	return Configuration{
		Port:   "8080",
		Secret: "dev-secret",
		URL:    "http://localhost:8080",
		FFmpeg: "/usr/local/bin/ffmpeg",
	}
}
