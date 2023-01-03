package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	ConfigServer     string
	ConfigServerPort string
}

func LoadConfig() *Config {
	env := os.Getenv("ACTIVE_ENV")
	if env == "" {
		env = "dev"
	}

	switch env {
	case "dev":
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
	case "prod":
	default:
		log.Fatal("Unknown environment")
	}

	log.Println("Loaded environment variables for", env)

	port := "6002"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	return &Config{
		ConfigServer:     os.Getenv("CONFIG_SERVER"),
		ConfigServerPort: os.Getenv("CONFIG_SERVER_PORT"),
		Port:             port,
	}
}
