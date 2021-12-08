package main

import (
	"context"
	_ "database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/go-redis/redis/v8"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jt-rose/clean_blog_server/graph"
	"github.com/jt-rose/clean_blog_server/graph/generated"

	// local imports
	initDB "github.com/jt-rose/clean_blog_server/database"
)

const defaultServerPort = "8080"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = defaultServerPort
	}

	dbpool := initDB.InitDB()
	defer dbpool.Close()

	err = dbpool.Ping(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Error: Redis connection failed")
	} else {
		fmt.Println(pong + " Redis connected")
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Redis connection failed: %v\n", err)
		os.Exit(1)
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", serverPort)
	log.Fatal(http.ListenAndServe(":"+serverPort, nil))
}
