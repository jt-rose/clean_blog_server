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
	CommentByUserID PaginatedCommentsLoader
	CommentByPostID PaginatedCommentsLoader
	CommentByCommentID PaginatedCommentsLoader // used for getting subcomments
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

func LoadCommentsByCommentID(ctx context.Context) func(ids []int) ([]model.PaginatedComments, []error){
	return func(ids []int) ([]model.PaginatedComments, []error) {
		// format user ids as SQL string param
		queryParam := utils.FormatSliceForSQLParams(ids)

		// attempt to fetch subcomments for each comment id
		comments, err := sql_models.Comments(qm.Where("response_to_comment_id = ANY(?::int[])", queryParam)).All(ctx, database.DB)
		if err != nil {
			return nil, []error{err}
		}

		// get the comment id for each comment found
		var currentCommentIDList []int
	for _, value := range comments {
		currentCommentIDList = append(currentCommentIDList, value.CommentID)
	}

	// format comment ids as SQL string param
	queryParam = utils.FormatSliceForSQLParams(currentCommentIDList)
	// check if subcomments exist that reference the current subcomment ids as a parent
	subComments, err := sql_models.Comments(qm.Where("response_to_comment_id = ANY(?::int[])", queryParam)).All(ctx, database.DB)
	if err != nil {
		return nil, []error{err}
	}

	// format comments for graphQL response and note if they have subcomments
	var formattedComments []model.Comment
	for _, comment := range comments {
		// check if comment has subcomments
		hasSubComment := false
		for _, subComment := range subComments {
			if subComment.ResponseToCommentID.Int == comment.CommentID {
				hasSubComment = true
			}
		}
		fmtComment := utils.ConvertComment(comment, hasSubComment)
		formattedComments = append(formattedComments, fmtComment)
	}
	
	// sort formatted comments by post ids from field resolver args
	// and format each association as a PaginatedComments response
	var formattedPaginatedComments []model.PaginatedComments
	for _, commentID := range ids {
		var commentsByCommentID []*model.Comment
		for _, comment := range formattedComments {
			if *comment.ResponseToCommentID == commentID {
				// preserve the current comments data before comment is updated
				currentComment := *&comment
				commentsByCommentID = append(commentsByCommentID, &currentComment)
			}
			
		}

		paginatedResponse := model.PaginatedComments{
			Comments: commentsByCommentID,
			More: false,
		}
		formattedPaginatedComments = append(formattedPaginatedComments, paginatedResponse)
	}

	return formattedPaginatedComments, nil
}
}

func LoadCommentsByPostID(ctx context.Context) func(ids []int) ([]model.PaginatedComments, []error){
	return func(ids []int) ([]model.PaginatedComments, []error) {
		// format post ids as SQL string param
		queryParam := utils.FormatSliceForSQLParams(ids)

		// attempt to fetch all top-level comments on each post
		comments, err := sql_models.Comments(qm.Where("post_id = ANY(?::int[]) AND response_to_comment_id IS NULL", queryParam)).All(ctx, database.DB)
		if err != nil {
			return nil, []error{err}
		}

		// get the comment id for each comment found
		var currentCommentIDList []int
	for _, value := range comments {
		currentCommentIDList = append(currentCommentIDList, value.CommentID)
	}

	// format comment ids as SQL string param
	queryParam = utils.FormatSliceForSQLParams(currentCommentIDList)
	// check if subcomments exist that reference the current comment ids as a parent
	subComments, err := sql_models.Comments(qm.Where("response_to_comment_id = ANY(?::int[])", queryParam)).All(ctx, database.DB)
	if err != nil {
		return nil, []error{err}
	}

	// format comments for graphQL response and note if they have subcomments
	var formattedComments []model.Comment
	for _, comment := range comments {
		// check if comment has subcomments
		hasSubComment := false
		for _, subComment := range subComments {
			if subComment.ResponseToCommentID.Int == comment.CommentID {
				hasSubComment = true
			}
		}
		fmtComment := utils.ConvertComment(comment, hasSubComment)
		formattedComments = append(formattedComments, fmtComment)
	}
	
	// sort formatted comments by post ids from field resolver args
	// and format each association as a PaginatedComments response
	var formattedPaginatedComments []model.PaginatedComments
	for _, postID := range ids {
		var commentsByPostID []*model.Comment
		for _, comment := range formattedComments {
			if comment.PostID == postID {
				// preserve the current comments data before comment is updated
				currentComment := *&comment
				commentsByPostID = append(commentsByPostID, &currentComment)
			}
			
		}

		paginatedResponse := model.PaginatedComments{
			Comments: commentsByPostID,
			More: false,
		}
		formattedPaginatedComments = append(formattedPaginatedComments, paginatedResponse)
	}

	return formattedPaginatedComments, nil
}
}

func LoadCommentsByUserID(ctx context.Context) func(ids []int) ([]model.PaginatedComments, []error){
	return func(ids []int) ([]model.PaginatedComments, []error) {
		// format user ids as SQL string param
		queryParam := utils.FormatSliceForSQLParams(ids)

		// attempt to fetch user comments
		comments, err := sql_models.Comments(qm.Where("user_id = ANY(?::int[]) AND deleted = false", queryParam)).All(ctx, database.DB)
		if err != nil {
			return nil, []error{err}
		}

		// get the comment id for each comment found
		var currentCommentIDList []int
	for _, value := range comments {
		currentCommentIDList = append(currentCommentIDList, value.CommentID)
	}

	// format comment ids as SQL string param
	queryParam = utils.FormatSliceForSQLParams(currentCommentIDList)
	// check if subcomments exist that reference the current comment ids as a parent
	subComments, err := sql_models.Comments(qm.Where("response_to_comment_id = ANY(?::int[])", queryParam)).All(ctx, database.DB)
	if err != nil {
		return nil, []error{err}
	}

	// format comments for graphQL response and note if they have subcomments
	var formattedComments []model.Comment
	for _, comment := range comments {
		// check if comment has subcomments
		hasSubComment := false
		for _, subComment := range subComments {
			if subComment.ResponseToCommentID.Int == comment.CommentID {
				hasSubComment = true
			}
		}
		fmtComment := utils.ConvertComment(comment, hasSubComment)
		formattedComments = append(formattedComments, fmtComment)
	}
	
	// sort formatted comments by post ids from field resolver args
	// and format each association as a PaginatedComments response
	var formattedPaginatedComments []model.PaginatedComments
	for _, userID := range ids {
		var commentsByUserID []*model.Comment
		for _, comment := range formattedComments {
			if comment.UserID == userID {
				// preserve the current comments data before comment is updated
				currentComment := *&comment
				commentsByUserID = append(commentsByUserID, &currentComment)
			}
			
		}

		paginatedResponse := model.PaginatedComments{
			Comments: commentsByUserID,
			More: false,
		}
		formattedPaginatedComments = append(formattedPaginatedComments, paginatedResponse)
	}

	return formattedPaginatedComments, nil
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
			CommentByUserID: PaginatedCommentsLoader{
				maxBatch: 100,
				wait:     1 * time.Millisecond,
				fetch: LoadCommentsByUserID(ginContext),
			},
			CommentByPostID: PaginatedCommentsLoader{
				maxBatch: 100,
				wait:     1 * time.Millisecond,
				fetch: LoadCommentsByPostID(ginContext),
			},
			CommentByCommentID: PaginatedCommentsLoader{
				maxBatch: 100,
				wait:     1 * time.Millisecond,
				fetch: LoadCommentsByCommentID(ginContext),
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
