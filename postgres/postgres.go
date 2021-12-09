package initDB

import (
	"context"
	_ "database/sql"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

const defaultDatabasePort = "5432"

func initDB() *pgxpool.Pool {
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