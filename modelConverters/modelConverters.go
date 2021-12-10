package modelConverters

import (
	gql_models "github.com/jt-rose/clean_blog_server/graph/model"
	sql_models "github.com/jt-rose/clean_blog_server/sql_models"
)

func ConvertPost(sql_post *sql_models.Post) gql_models.Post {
	return gql_models.Post{
		UserID: sql_post.UserID,
		PostID: sql_post.PostID,
		Title: sql_post.Title,
		Subtitle: sql_post.Subtitle,
		PostText: sql_post.PostText,
		CreatedAt: sql_post.CreatedAt,
	}
}
/*
func ConvertComment(sql_comment *sql_models.Comment) gql_models.Comment {
	return gql_models.Comment{
		CommentID: sql_comment.CommentID,
		ResponseToCommentID: sql_comment.ResponseToCommentID, // change to 0's, nonnullable
		PostID: sql_comment.PostID,
		UserID: sql_comment.UserID,
		CommentText: sql_comment.CommentText,
		CreatedAt: sql_comment.CreatedAt,
	}
}*/

func ConvertUser(sql_user *sql_models.User) gql_models.User {
	return gql_models.User{
		UserID: sql_user.UserID,
		Username: sql_user.Username,
		Email: sql_user.Email,
		CreatedAt: sql_user.CreatedAt,
	}
}

func ConvertCommentVote() {}

func ConvertPostVote() {}