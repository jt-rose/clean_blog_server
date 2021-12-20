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
	SESSION_KEY   string
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
		SESSION_KEY:      os.Getenv("SESSION_KEY"),
	}

	if ENV_VAR.DATABASE_URL == "" || ENV_VAR.DATABASE_PORT == "" || ENV_VAR.SERVER_PORT == "" || ENV_VAR.SESSION_KEY == "" {
		log.Fatal("Error loading .env file")
	}

	return ENV_VAR
}

var ENV_VARIABLES ENV_Variables = loadEnvVariables()
var COOKIE_NAME="cid"
var UNAUTHENTICATED_ERROR_MESSAGE = "Must be logged in!"
var ONLY_AUTHOR_ALLOWED_ERROR_MESSAGE = "Only the author of the blog may add, edit, or delete posts"
var ONLY_COMMENT_AUTHOR_MAY_EDIT = "Only the author of a comment can edit, delete, or restore it"
// May refactor later to use a less "hard-coded" solution
var AUTHOR_ID = 2
