package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	//"time"

	"github.com/joho/godotenv"
	"github.com/jt-rose/clean_blog_server/graph/generated"
	"github.com/jt-rose/clean_blog_server/graph/model"
	convert "github.com/jt-rose/clean_blog_server/modelConverters"

	sql_models "github.com/jt-rose/clean_blog_server/sql_models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// Open handle to database like normal
  func initDB() *sql.DB {
	err := godotenv.Load()
	if err != nil {
	  log.Fatal("Error loading .env file")
	}
  
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
	fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	  os.Exit(1)
	}
  
	return db
  }
  
  var DB = initDB()

  func handleSQLErrors(ctx context.Context, err error, errOrigin string) error {
		if err.Error() == "sql: no rows in result set" {
			return errors.New("no posts found")
		} else {
			fmt.Println("Error: ", err)
			errorLog := sql_models.ErrorLog{
				ErrMessage: err.Error(),
				ErrorOrigin: errOrigin,
			}
			errorLog.Insert(ctx, DB, boil.Infer())
			return errors.New("data unavailable")
		}
  }

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
	post, err := sql_models.Posts(qm.Where("post_id = ?", postID)).One(ctx, DB)

  if post == nil {
	return nil, handleSQLErrors(ctx, err, "GetPost")
  }

  formattedPost := convert.ConvertPost(post)
  return &formattedPost, nil
}

func (r *queryResolver) GetUser(ctx context.Context, userID int) (*model.User, error) {
	user, err := sql_models.Users(qm.Where("user_id = ?", userID)).One(ctx, DB)

  if user == nil {
	return nil, handleSQLErrors(ctx, err, "GetUser")
  }

  formattedUser := convert.ConvertUser(user)
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
