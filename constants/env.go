package constants

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ENV_Variables struct {
	DATABASE_URL  string
	DATABASE_PORT string
	SERVER_PORT   string
	FRONTEND_URL string
	SESSION_KEY   string
	EMAIL_ADDRESS string
	EMAIL_PASSWORD string
}

func loadEnvVariables() ENV_Variables {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ENV_VAR := ENV_Variables{
		DATABASE_URL:  os.Getenv("DATABASE_URL"),
		DATABASE_PORT: os.Getenv("DATABASE_URL"),
		SERVER_PORT:   os.Getenv("SERVER_PORT"),
		FRONTEND_URL: os.Getenv("FRONTEND_URL"),
		SESSION_KEY:      os.Getenv("SESSION_KEY"),
		EMAIL_ADDRESS: os.Getenv("EMAIL_ADDRESS"),
		EMAIL_PASSWORD: os.Getenv("EMAIL_PASSWORD"),
	}

	if ENV_VAR.DATABASE_URL == "" || ENV_VAR.DATABASE_PORT == "" || ENV_VAR.SERVER_PORT == "" || ENV_VAR.SESSION_KEY == "" {
		log.Fatal("Error loading .env file")
	}

	return ENV_VAR
}

var ENV_VARIABLES ENV_Variables = loadEnvVariables()
var COOKIE_NAME="cid"