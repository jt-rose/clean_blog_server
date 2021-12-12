package errorHandler

import (
	"context"
	"errors"
	"fmt"
	"runtime"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"

	postgres "github.com/jt-rose/clean_blog_server/postgres"
	sql_models "github.com/jt-rose/clean_blog_server/sql_models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func HandleErrors(ctx context.Context, e error) *gqlerror.Error {
	err := graphql.DefaultErrorPresenter(ctx, e)

	//var myErr *MyError
	if errors.Is(e, errors.New("sql: no rows in result set")) {
		err.Message = "No matching data found in database"
	} else {
		// provide generic response to hide error details from the client
		err.Message = "data unavailable"

		// print out data on point of failure
		pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	fmt.Printf("Error Found: %s:%d %s\n", frame.File, frame.Line, frame.Function)

	// add data on point of failure to error log
	errorLog := sql_models.ErrorLog{
		ErrMessage: err.Error(),
		ErrorOrigin: frame.Function,
	}
	errorLog.Insert(ctx, postgres.DB, boil.Infer())
	}

	// return newly formatted error
	return err
}
