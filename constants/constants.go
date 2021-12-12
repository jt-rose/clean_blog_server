package constants

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ENV_Variables struct {
	DATABASE_URL string
	DATABASE_PORT string
	SERVER_PORT string
}

func loadEnvVariables() ENV_Variables {
	err := godotenv.Load()
	if err != nil {
	  log.Fatal("Error loading .env file")
	}

	ENV_VAR := ENV_Variables{
		DATABASE_URL: os.Getenv("DATABASE_URL"),
	DATABASE_PORT: os.Getenv("DATABASE_URL"),
	SERVER_PORT: os.Getenv("SERVER_PORT"),
	}

	if ENV_VAR.DATABASE_URL == "" || ENV_VAR.DATABASE_PORT == "" || ENV_VAR.SERVER_PORT == "" {
		log.Fatal("Error loading .env file")
	}

	return ENV_VAR
}

var ENV_VARIABLES ENV_Variables = loadEnvVariables()