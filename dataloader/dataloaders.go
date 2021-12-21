package dataloader

import (
	"context"
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
	CommentByUserID CommentLoader
	CommentByPostID CommentLoader
	VotesByPostID VotesLoader
	VotesByCommentID VotesLoader
}

func LoadUsers(ctx context.Context) func(ids []int) ([]model.User, []error){
	return func(ids []int) ([]model.User, []error) {
		// format ids as SQL string param
		queryParam := utils.FormatSliceForSQLParams(ids)

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

func LoadVotesByCommentID(ctx context.Context) func(ids []int) ([]model.Votes, []error){
	return func(ids []int) ([]model.Votes, []error) {
		// format ids as SQL string param
		queryParam := utils.FormatSliceForSQLParams(ids)

		// attempt to fetch users
		commentVotes, err := sql_models.CommentVotes(qm.Where("comment_id = ANY(?::int[])", queryParam)).All(ctx, database.DB)
		
		// format error
		formattedErrors := []error{err}
		if err != nil {
			return nil, formattedErrors
		}
		
		// total upvote and downvote counts by comment_id
		voteCounts := make([]model.Votes, len(ids))
		for i, commentID := range ids {
			votes := model.Votes{}
			for _, singleVote := range commentVotes {
				if singleVote.CommentID == commentID {
					if singleVote.VoteValue == -1 {
						votes.Downvote++
					}
					if singleVote.VoteValue == 1 {
						votes.Upvote++
					}
				}
			}
			voteCounts[i] = votes
		}

		// return vote counts
		return voteCounts, nil
}
}

func LoadVotesByPostID(ctx context.Context) func(ids []int) ([]model.Votes, []error){
	return func(ids []int) ([]model.Votes, []error) {
		// format ids as SQL string param
		queryParam := utils.FormatSliceForSQLParams(ids)

		// attempt to fetch users
		postVotes, err := sql_models.PostVotes(qm.Where("post_id = ANY(?::int[])", queryParam)).All(ctx, database.DB)
		
		// format error
		formattedErrors := []error{err}
		if err != nil {
			return nil, formattedErrors
		}
		
		// total upvote and downvote counts by comment_id
		voteCounts := make([]model.Votes, len(ids))
		for i, postID := range ids {
			votes := model.Votes{}
			for _, singleVote := range postVotes {
				if singleVote.PostID == postID {
					if singleVote.VoteValue == -1 {
						votes.Downvote++
					}
					if singleVote.VoteValue == 1 {
						votes.Upvote++
					}
				}
			}
			voteCounts[i] = votes
		}

		// return vote counts
		return voteCounts, nil
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
			VotesByCommentID: VotesLoader{
				maxBatch: 100,
				wait:     1 * time.Millisecond,
				fetch: LoadVotesByCommentID(ginContext),
			},
			VotesByPostID: VotesLoader{
				maxBatch: 100,
				wait:     1 * time.Millisecond,
				fetch: LoadVotesByPostID(ginContext),
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
