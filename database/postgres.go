package database

import (
	"database/sql"
	"fmt"
	"os"

	ENV "github.com/jt-rose/clean_blog_server/constants"
	_ "github.com/lib/pq"
)

const defaultDatabasePort = "5432"

func initDB() *sql.DB {

	db, err := sql.Open("postgres", ENV.ENV_VARIABLES.DATABASE_URL)
	if err != nil {
	fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	  os.Exit(1)
	} else {
		fmt.Println("Connected to Postgres database on port " + ENV.ENV_VARIABLES.DATABASE_PORT)
	}
	err = db.Ping()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	  os.Exit(1)
	}
	return db
  }
  
  var DB = initDB()