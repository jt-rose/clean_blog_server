package main

import (
	// graphQL handlers
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/playground"

	// gqlgen generated models
	"github.com/jt-rose/clean_blog_server/graph"
	"github.com/jt-rose/clean_blog_server/graph/generated"

	// gin + redis session middleware
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"

	// local imports
	helmet "github.com/danielkov/gin-helmet"
	auth "github.com/jt-rose/clean_blog_server/auth"
	ENV "github.com/jt-rose/clean_blog_server/constants"
	errorHandler "github.com/jt-rose/clean_blog_server/errorHandler"
	postgres "github.com/jt-rose/clean_blog_server/postgres"
)

// Defining the Graphql handler
func graphqlHandler() gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {

	DB := postgres.DB
	defer DB.Close()

	// Setting up Gin
	r := gin.Default()
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte(ENV.ENV_VARIABLES.SESSION_KEY))
	
	
	r.Use(sessions.Sessions("session_id", store))
	r.Use(auth.GinContextToContextMiddleware())
	r.Use(auth.Authenticate())
	r.Use(helmet.Default())
	
	r.GET("/inc", auth.TESTREDIS)
	r.POST("/query", graphqlHandler())
	r.GET("/", playgroundHandler())

	// initialize GraphQL server
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	// set up error and panic handling
	srv.SetErrorPresenter(errorHandler.HandleErrors)
	srv.SetRecoverFunc(errorHandler.HandlePanics)
	// limit query complexity to depth of 20
	srv.Use(extension.FixedComplexityLimit(20))

	r.Run()
}
