package initDB

import (
	"context"
	_ "database/sql"
	"fmt"
	"os"

	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

const defaultDatabasePort = "5432"

func initDB() *pgxpool.Pool {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	databasePort := os.Getenv("DB_PORT")
	if databasePort == "" {
		databasePort = defaultDatabasePort
	}

	dbpool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Println("Postgres database connected on port " + databasePort)
	}

	return dbpool
}

var DBPool = initDB()