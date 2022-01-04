package middleware

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
)

func HandleLogs(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	oc := graphql.GetOperationContext(ctx)
	fmt.Printf("[GQL] %s: ", oc.OperationName)
	return next(ctx)
}