package middleware

import (
	"context"
	"fmt"

	sessions "github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jt-rose/clean_blog_server/graph/model"
)

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

// Authenticate decodes the share session cookie and packs the session into context

func Authenticate() gin.HandlerFunc {
	return func(ginContext *gin.Context) {

		// connext to existing session or generate a new one
		// a session will always be returned
		session := sessions.Default(ginContext)

		// Retrieve our User struct and type-assert it
		val := session.Get("user")
		/*var user = &gql_models.User{}
		user, ok := val.(*gql_models.User)

		// if the struct type does not match the expected User struct
		// invalidate the cookie, store err in error log,
		// and move on without adding to context
		if !ok {
			session.Delete("user")
			err := session.Save()
			// TODO: add to error log
			if err != nil {
				http.Error(ginContext.Writer, err.Error(), http.StatusInternalServerError)
			}
			ginContext.Next()
			return
		}*/

		// pass user infos to context
		ctx := context.WithValue(ginContext.Request.Context(), userCtxKey, val)

		// and call next with our new context
		ginContext.Request = ginContext.Request.WithContext(ctx)
		ginContext.Next()
	}
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *model.User {
	raw, _ := ctx.Value(userCtxKey).(*model.User)
	return raw
}

// add gin context
func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "GinContextKey", c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value("GinContextKey")
	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return gc, nil
}

// return gin context and sessions
// handle error internally to reduce boilerplate
func GetGinContextAndSessions(ctx context.Context) (*gin.Context, sessions.Session) {
	// get gin context
	gc, err := GinContextFromContext(ctx)
	if err != nil {
return nil, nil
	}
	// get sessions
	session := sessions.Default(gc)

	return gc,
		session
}