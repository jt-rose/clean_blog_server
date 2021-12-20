package main

import (
	// graphQL handlers
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/playground"

	// gqlgen generated models

	"github.com/jt-rose/clean_blog_server/dataloader"
	"github.com/jt-rose/clean_blog_server/graph"
	"github.com/jt-rose/clean_blog_server/graph/generated"

	// gin + redis session middleware
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"

	// local imports
	helmet "github.com/danielkov/gin-helmet"
	ENV "github.com/jt-rose/clean_blog_server/constants"
	database "github.com/jt-rose/clean_blog_server/database"
	middleware "github.com/jt-rose/clean_blog_server/middleware"
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
	// get Database and close on defer
	DB := database.DB
	defer DB.Close()

	// Setting up Gin
	r := gin.Default()
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte(ENV.ENV_VARIABLES.SESSION_KEY))
	
	// set up middleware
	r.Use(sessions.Sessions("session_id", store))
	r.Use(middleware.GinContextToContextMiddleware())
	r.Use(middleware.Authenticate())
	r.Use(dataloader.UseDataLoaders())
	r.Use(helmet.Default())
	
	// set up routes
	r.POST("/query", graphqlHandler())
	r.GET("/", playgroundHandler())

	// initialize GraphQL server
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	
	// set up error and panic handling
	srv.SetErrorPresenter(middleware.HandleErrors)
	srv.SetRecoverFunc(middleware.HandlePanics)
	
	// limit query complexity to depth of 20
	srv.Use(extension.FixedComplexityLimit(20))

	// run on default available ports
	r.Run()
}
