package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jt-rose/clean_blog_server/graph"
	"github.com/jt-rose/clean_blog_server/graph/generated"

	// local imports
	ENV "github.com/jt-rose/clean_blog_server/constants"
	postgres "github.com/jt-rose/clean_blog_server/postgres"
	redis "github.com/jt-rose/clean_blog_server/redis"
)

func main() {

	DB := postgres.DB
	defer DB.Close()

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
	//srv.SetErrorPresenter()

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", ENV.ENV_VARIABLES.SERVER_PORT)
	log.Fatal(http.ListenAndServe(":"+ENV.ENV_VARIABLES.SERVER_PORT, nil))
}
