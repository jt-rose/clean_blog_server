package middleware

import (
	"context"
	"errors"
	"fmt"
	"runtime"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/jt-rose/clean_blog_server/constants"
	database "github.com/jt-rose/clean_blog_server/database"
	sql_models "github.com/jt-rose/clean_blog_server/sql_models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// read the function environment and store a SQL record of the error
func storeErrorLog(ctx context.Context, err error) error {
	// print out data on point of failure
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	fmt.Println("Error Encountered: ", err.Error())
	fmt.Printf("Error Found at: %s:%d %s\n", frame.File, frame.Line, frame.Function)

	// add data on point of failure to error log
	errorLog := sql_models.ErrorLog{
		ErrMessage:  err.Error(),
		ErrorOrigin: frame.Function,
	}
	errorLog.Insert(ctx, database.DB, boil.Infer())
	return err
}

// parse errors and hide select error messages from end user
func HandleErrors(ctx context.Context, e error) *gqlerror.Error {
	err := graphql.DefaultErrorPresenter(ctx, e)

	// if data not found, return user-friendly error message
	if err.Message == "sql: no rows in result set" {
		err.Message = "No matching data found in database"
		// if custom error not found, store the issue in the error log
		// and hide the details of the error from the client
	} else if !constants.IsCustomError(err.Message) {
		storeErrorLog(ctx, err)
		// provide generic response to hide error details from the client
		err.Message = "data currently unavailable"

	}
	// return newly formatted error or custom error
	return err
}

// store error log and format message when recovering from a panic
// to be used with the gql server.SetRecoverFunc function
func HandlePanics(ctx context.Context, err interface{}) error {
	var foundError error
	genericError := errors.New("Internal server error!")
	// run type assertion to confirm err is a map
	errorStruct, ok := err.(error)
	if ok {
		foundError = errorStruct
	} else {
		foundError = genericError
	}

	// notify bug tracker and print to console
	// store detailed error message for error log
	// but only show "Internal Server Error" for end users
	fmt.Println(foundError.Error())
	storeErrorLog(ctx, foundError)

	return genericError
}
