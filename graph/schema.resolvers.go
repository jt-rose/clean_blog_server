package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"

	"github.com/jt-rose/clean_blog_server/graph/generated"
	"github.com/jt-rose/clean_blog_server/graph/model"
	models "github.com/jt-rose/clean_blog_server/sql_models"
	utils "github.com/jt-rose/clean_blog_server/utils"

	sessions "github.com/gin-contrib/sessions"
	database "github.com/jt-rose/clean_blog_server/database"
	middleware "github.com/jt-rose/clean_blog_server/middleware"
	sql_models "github.com/jt-rose/clean_blog_server/sql_models"
	boil "github.com/volatiletech/sqlboiler/v4/boil"
	qm "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (r *commentResolver) User(ctx context.Context, obj *model.Comment) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *commentResolver) Votes(ctx context.Context, obj *model.Comment) (*model.Votes, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddPost(ctx context.Context, postInput model.PostInput) (*model.Post, error) {
	/*gc, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}*/
	// add validation
	newPost := models.Post{
		UserID: 1,///////// add context/userID
		Title: postInput.Title,
		Subtitle: *postInput.Subtitle,
		PostText: postInput.Text,
	}
	err := newPost.Insert(ctx, database.DB, boil.Infer())
	if err != nil {
		return nil, err
	}

	formattedPost := utils.ConvertPost(&newPost)
	return &formattedPost, nil
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
	gc, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}
	
	hashedPassword, err := utils.HashPassword(userInput.Password)

	if err != nil {
		// TODO: add error log
		return nil, errors.New("Password error: please contact administrator")
	}
	
	newUser := sql_models.User{
		Username: userInput.Username,
		Email: userInput.Email,
		UserPassword: hashedPassword,
	}

	err = newUser.Insert(ctx, database.DB, boil.Infer())

	if err != nil {
		// TODO: add error log and handling
		return nil, err
	}
	

	// format user and remove password from struct
	formattedUser := utils.ConvertUser(&newUser)

	// get session and add new user 
	// add err handling
	session := sessions.Default(gc)
	session.Set("user", formattedUser.UserID)
	err = session.Save()
	if err != nil {
		return nil, err
	}
	// secure and add to redis session

	return &formattedUser, err
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
	post, err := sql_models.Posts(qm.Where("post_id = ?", postID)).One(ctx, database.DB)

  if post == nil {
	return nil, err
  }

  formattedPost := utils.ConvertPost(post)
  return &formattedPost, nil
}

func (r *queryResolver) GetUser(ctx context.Context, userID int) (*model.User, error) {
	user, err := sql_models.Users(qm.Where("user_id = ?", userID)).One(ctx, database.DB)

  if user == nil {
	return nil, err
  }

  formattedUser := utils.ConvertUser(user)
  return &formattedUser, nil
}

func (r *queryResolver) GetManyPosts(ctx context.Context, postSearch model.PostSearch) (*model.PaginatedPosts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetManyUsers(ctx context.Context, userSearch model.UserSearch) (*model.PaginatedUsers, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetManyComments(ctx context.Context, commentSearch model.CommentSearch) (*model.PaginatedComments, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Me(ctx context.Context, userID int) (bool, error) {
	gc, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return false, err
	}

	// format user and remove password from struct
	//formattedUser := modelConverters.ConvertUser(&newUser)
	session := sessions.Default(gc)
	user := session.Get("user")
	if user == "" || user == nil {
		return false, nil
	}
 	return true, nil
	//session.Save(gc.Request, gc.Writer)
	//panic(fmt.Errorf("not implemented"))
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
