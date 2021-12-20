package dataloader

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	database "github.com/jt-rose/clean_blog_server/database"
	"github.com/jt-rose/clean_blog_server/graph/model"
	sql_models "github.com/jt-rose/clean_blog_server/sql_models"
	utils "github.com/jt-rose/clean_blog_server/utils"
	qm "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

const loadersKey = "dataloaders"

type Loaders struct {
	UserById UserLoader
}

func LoadUsers(ctx context.Context) func(ids []int) ([]model.User, []error){
	return func(ids []int) ([]model.User, []error) {
		// convert ids to []string
		var stringArgs []string
		for _, id := range ids {
			stringArgs = append(stringArgs, strconv.Itoa(id))
		}

		// format param of SQL query
		queryParam := "{"+ strings.Join(stringArgs, ",") + "}"

		// attempt to fetch users
		users, err := sql_models.Users(qm.Where("user_id = ANY(?::int[])", queryParam)).All(ctx, database.DB)
		
		// format users and error
		var formattedUsers []model.User
		formattedErrors := []error{err}
		
		if err != nil {
			return nil, formattedErrors
		}
		
		for _, user := range users {
			fmtUser := utils.ConvertUser(user)
			formattedUsers = append(formattedUsers, fmtUser)
		}

		// sort users according to order of ids in dataloader arg
		sortedUsers := make([]model.User, len(ids))
		for i, id := range ids {
			for _, user := range formattedUsers {
				if user.UserID == id {
					sortedUsers[i] = user
				}
			}
			
		}
		
		// return sorted users
		return sortedUsers, nil
}
}

func UseDataLoaders() gin.HandlerFunc {
	return func(ginContext *gin.Context) {

		// pass dataloaders to context
		ctx := context.WithValue(ginContext.Request.Context(), loadersKey, &Loaders{
			UserById: UserLoader{
				maxBatch: 100,
				wait:     1 * time.Millisecond,
				fetch: LoadUsers(ginContext),
			},
		})

		// and call next with our new context
		ginContext.Request = ginContext.Request.WithContext(ctx)
		ginContext.Next()
	}
}

func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}
