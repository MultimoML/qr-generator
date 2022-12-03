package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Environment() {
	switch os.Getenv("ENVIRONMENT") {
	case "dev":
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading environment file")
		} else {
			log.Println("Loaded dev environment")
		}
	case "prod":
		log.Println("Loaded prod environment")
	default:
		log.Fatal("Environment variable ENVIRONMENT is not set or set to wrong value")
	}
}
