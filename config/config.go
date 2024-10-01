package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Configuration struct {
	Port   int
	Secret string
	URL    string
	FFmpeg string
}

var Config Configuration

func init() {
	godotenv.Load()

	nodeEnv := os.Getenv("NODE_ENV")
	if nodeEnv == "" {
		nodeEnv = "development"
	}

	switch nodeEnv {
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
		Port:   8080,
		Secret: "dev-secret",
		URL:    "http://localhost:8080",
		FFmpeg: "/usr/local/bin/ffmpeg",
	}
}

func getStagingConfig() Configuration {
	// Implement staging configuration
	return Configuration{
		// ... staging config values
	}
}

func getProductionConfig() Configuration {
	// Implement production configuration
	return Configuration{
		// ... production config values
	}
}
