package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jt-rose/clean_blog_server/graph"
	"github.com/jt-rose/clean_blog_server/graph/generated"

	// local imports
	postgres "github.com/jt-rose/clean_blog_server/postgres"
	redis "github.com/jt-rose/clean_blog_server/redis"
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

	//dbpool := initDB.InitDB()
	//defer dbpool.Close()

	// remove later
	err = postgres.DBPool.Ping(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	rdb := redis.RedisClient
	// remove later
	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Redis connection failed: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Println(pong + " Redis connected")
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", serverPort)
	log.Fatal(http.ListenAndServe(":"+serverPort, nil))
}
