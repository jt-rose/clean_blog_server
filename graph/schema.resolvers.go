package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/jt-rose/clean_blog_server/graph/generated"
	"github.com/jt-rose/clean_blog_server/graph/model"
)

func (r *commentResolver) User(ctx context.Context, obj *model.Comment) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *commentResolver) Votes(ctx context.Context, obj *model.Comment) (*model.Votes, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddPost(ctx context.Context, postInput model.PostInput) (*model.Post, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) EditPost(ctx context.Context, postID int, postInput model.PostInput) (*model.Post, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeletePost(ctx context.Context, postID int) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddComment(ctx context.Context, postID int, responseToCommentID *int, commentText string) (*model.Comment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) EditComment(ctx context.Context, commentID int, newCommentText string) (*model.Comment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteComment(ctx context.Context, commentID int) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VoteOnPost(ctx context.Context, postID int, voteValue int) (*model.PostVote, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VoteOnComment(ctx context.Context, commentID int, voteValue int) (*model.CommentVote, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RegisterNewUser(ctx context.Context, userInput model.UserInput) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Login(ctx context.Context, username string, password string) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Logout(ctx context.Context) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ForgotPassword(ctx context.Context, username string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ResetPassword(ctx context.Context, username string, newPassword string) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *postResolver) User(ctx context.Context, obj *model.Post) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *postResolver) Votes(ctx context.Context, obj *model.Post) (*model.Votes, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetPost(ctx context.Context, postID int) (*model.Post, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetManyPosts(ctx context.Context, postIds []int) ([]*model.Post, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetComments(ctx context.Context, postID int) ([]*model.Comment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Me(ctx context.Context, userID int) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) IsAuthor(ctx context.Context, userID int) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) Posts(ctx context.Context, obj *model.User) ([]*model.Post, error) {
	panic(fmt.Errorf("not implemented"))
}

// Comment returns generated.CommentResolver implementation.
func (r *Resolver) Comment() generated.CommentResolver { return &commentResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Post returns generated.PostResolver implementation.
func (r *Resolver) Post() generated.PostResolver { return &postResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type commentResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type postResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
