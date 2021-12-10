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

func ConvertComment() {}

func ConvertUser() {}

func ConvertCommentVote() {}

func ConvertPostVote() {}