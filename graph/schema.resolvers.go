package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"net/url"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jt-rose/clean_blog_server/constants"
	"github.com/jt-rose/clean_blog_server/graph/generated"
	"github.com/jt-rose/clean_blog_server/graph/model"
	utils "github.com/jt-rose/clean_blog_server/utils"

	database "github.com/jt-rose/clean_blog_server/database"
	dataloader "github.com/jt-rose/clean_blog_server/dataloader"
	middleware "github.com/jt-rose/clean_blog_server/middleware"
	sql_models "github.com/jt-rose/clean_blog_server/sql_models"
	null "github.com/volatiletech/null/v8"
	boil "github.com/volatiletech/sqlboiler/v4/boil"
	qm "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

/* -------------------------------------------------------------------------- */
/*                         dataloader field resolvers                         */
/* -------------------------------------------------------------------------- */

func (r *commentResolver) User(ctx context.Context, obj *model.Comment) (*model.User, error) {
	user, err := dataloader.For(ctx).UserById.Load(obj.UserID)
	return &user, err
}

func (r *commentResolver) Comments(ctx context.Context, obj *model.Comment) (*model.PaginatedComments, error) {
	paginatedComments, err := dataloader.For(ctx).CommentByCommentID.Load(obj.CommentID)
	return &paginatedComments, err
}

func (r *commentResolver) Votes(ctx context.Context, obj *model.Comment) (*model.Votes, error) {
	votes, err := dataloader.For(ctx).VotesByCommentID.Load(obj.CommentID)
	return &votes, err
}

/* -------------------------------------------------------------------------- */
/*                                  Post CRUD                                 */
/* -------------------------------------------------------------------------- */

func (r *mutationResolver) AddPost(ctx context.Context, postInput model.PostInput, authorID int) (*model.Post, error) {
	// confirm user is the author of the blog
	isAuthor, userID, err := middleware.ConfirmAuthor(ctx, authorID)
	if err != nil {
		return nil, err
	}

	if !isAuthor {
		return nil, errors.New(constants.ONLY_AUTHOR_ALLOWED_ERROR_MESSAGE)
	}

	// attenpt to add new post
	newPost := sql_models.Post{
		UserID:    userID,
		Title:     postInput.Title,
		PostText:  postInput.Text,
		Published: postInput.Published,
	}

	if postInput.Subtitle != nil {
		newPost.Subtitle = *postInput.Subtitle
	}

	err = newPost.Insert(ctx, database.DB, boil.Infer())
	if err != nil {
		return nil, err
	}

	// return gql version of the post
	gql_post := utils.ConvertPost(&newPost)
	return &gql_post, nil
}

func (r *mutationResolver) EditPost(ctx context.Context, postID int, postInput model.PostInput, authorID int) (*model.Post, error) {
	// confirm user is the author of the blog
	isAuthor, _, err := middleware.ConfirmAuthor(ctx, authorID)
	if err != nil {
		return nil, err
	}

	if !isAuthor {
		return nil, errors.New(constants.ONLY_AUTHOR_ALLOWED_ERROR_MESSAGE)
	}

	// attempt to update the post in the database
	currentPost, err := sql_models.FindPost(ctx, database.DB, postID)
	if err != nil {
		return nil, err
	}

	currentPost.Title = postInput.Title
	currentPost.Subtitle = *postInput.Subtitle
	currentPost.PostText = postInput.Text
	currentPost.Published = postInput.Published

	_, err = currentPost.Update(ctx, database.DB, boil.Infer())
	if err != nil {
		return nil, err
	}

	// return gql version of sql post object
	gql_post := utils.ConvertPost(currentPost)
	return &gql_post, nil
}

func (r *mutationResolver) DeletePost(ctx context.Context, postID int, authorID int) (bool, error) {
	// confirm user is the author of the blog
	isAuthor, _, err := middleware.ConfirmAuthor(ctx, authorID)
	if err != nil {
		return false, err
	}

	if !isAuthor {
		return false, errors.New(constants.ONLY_AUTHOR_ALLOWED_ERROR_MESSAGE)
	}

	// attempt to update the deleted property to true
	_, err = sql_models.Posts(qm.Where("post_id = ?", postID)).UpdateAll(ctx, database.DB, sql_models.M{"deleted": true})

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) RestorePost(ctx context.Context, postID int, authorID int) (bool, error) {
	// confirm user is the author of the blog
	isAuthor, _, err := middleware.ConfirmAuthor(ctx, authorID)
	if err != nil {
		return false, err
	}

	if !isAuthor {
		return false, errors.New(constants.ONLY_AUTHOR_ALLOWED_ERROR_MESSAGE)
	}

	// attempt to restore post by updating deleted property to false
	_, err = sql_models.Posts(qm.Where("post_id = ?", postID)).UpdateAll(ctx, database.DB, sql_models.M{"deleted": false})

	if err != nil {
		return false, err
	}

	return true, nil
}

/* -------------------------------------------------------------------------- */
/*                                comment CRUD                                */
/* -------------------------------------------------------------------------- */

func (r *mutationResolver) AddComment(ctx context.Context, postID int, responseToCommentID *int, commentText string) (*model.Comment, error) {
	// confirm authenticated
	userID, err := middleware.GetUserIDFromSessions(ctx)
	if err != nil {
		return nil, err
	}

	// attempt to add comment to database
	newComment := sql_models.Comment{
		UserID: userID,
		PostID: postID,
		ResponseToCommentID: null.Int{
			Int:   *responseToCommentID,
			Valid: *responseToCommentID == 0,
		},
		CommentText: commentText,
	}

	err = newComment.Insert(ctx, database.DB, boil.Infer())

	if err != nil {
		return nil, err
	}
	// return graphQL version of comment object
	gql_comment := utils.ConvertComment(&newComment, false)
	return &gql_comment, nil
}

func (r *mutationResolver) EditComment(ctx context.Context, commentID int, newCommentText string) (*model.Comment, error) {
	// authenticate user
	userID, err := middleware.GetUserIDFromSessions(ctx)
	if err != nil {
		return nil, err
	}

	// confirm user is author of comment
	comment, err := sql_models.Comments(qm.Where("comment_id = ?", commentID)).One(ctx, database.DB)
	if err != nil {
		return nil, err
	}

	// reject if not the author of the comment
	if comment.UserID != userID {
		err = errors.New(constants.ONLY_COMMENT_AUTHOR_MAY_EDIT)
		return nil, err
	}

	// attempt to edit the comment in the database
	_, err = sql_models.Comments(qm.Where("comment_id = ?", commentID)).UpdateAll(ctx, database.DB, sql_models.M{"comment_text": newCommentText})
	if err != nil {
		return nil, err
	}

	// return gql version of updated comment
	// there is the small
	hasSubComments, _ := sql_models.Comments(qm.Where("response_to_comment_id = ?", commentID)).Exists(ctx, database.DB)

	// if the above query has an error, we will simply use the zero value for hasSubComments
	gql_comment := utils.ConvertComment(comment, hasSubComments)
	gql_comment.CommentText = newCommentText
	return &gql_comment, nil
}

func (r *mutationResolver) DeleteComment(ctx context.Context, commentID int) (bool, error) {
	// authenticate user
	userID, err := middleware.GetUserIDFromSessions(ctx)
	if err != nil {
		return false, err
	}

	// confirm user is author of comment
	comment, err := sql_models.Comments(qm.Where("comment_id = ?", commentID)).One(ctx, database.DB)
	if err != nil {
		return false, err
	}

	// reject if not the author of the comment
	if comment.UserID != userID {
		err = errors.New(constants.ONLY_COMMENT_AUTHOR_MAY_EDIT)
		return false, err
	}

	// attempt to delete the comment in the database
	_, err = sql_models.Comments(qm.Where("comment_id = ?", commentID)).UpdateAll(ctx, database.DB, sql_models.M{"deleted": false})
	if err != nil {
		return false, err
	}

	// return boolean confirming successful restoration
	return true, nil
}

func (r *mutationResolver) RestoreComment(ctx context.Context, commentID int) (bool, error) {
	// authenticate user
	userID, err := middleware.GetUserIDFromSessions(ctx)
	if err != nil {
		return false, err
	}

	// confirm user is author of comment
	comment, err := sql_models.Comments(qm.Where("comment_id = ?", commentID)).One(ctx, database.DB)
	if err != nil {
		return false, err
	}

	// reject if not the author of the comment
	if comment.UserID != userID {
		err = errors.New(constants.ONLY_COMMENT_AUTHOR_MAY_EDIT)
		return false, err
	}

	// attempt to delete the comment in the database
	_, err = sql_models.Comments(qm.Where("comment_id = ?", commentID)).UpdateAll(ctx, database.DB, sql_models.M{"deleted": true})
	if err != nil {
		return false, err
	}

	// return boolean confirming successful deletion
	return true, nil
}

/* -------------------------------------------------------------------------- */
/*                             vote functionality                             */
/* -------------------------------------------------------------------------- */

func (r *mutationResolver) VoteOnPost(ctx context.Context, postID int, voteValue model.VoteValue) (*model.PostVote, error) {
	// authenticate user
	userID, err := middleware.GetUserIDFromSessions(ctx)
	if err != nil {
		return nil, err
	}

	// check if user has already voted
	currentPostVote, err := sql_models.FindPostVote(ctx, database.DB, postID, userID)

	if err != nil {
		return nil, err
	}

	// attempt to add new vote or update existing vote
	if currentPostVote == nil {
		currentPostVote = &sql_models.PostVote{
			PostID:    postID,
			UserID:    userID,
			VoteValue: utils.ConvertGQLVoteValueEnums(voteValue),
		}
		err = currentPostVote.Insert(ctx, database.DB, boil.Infer())
		if err != nil {
			return nil, err
		}

	} else {
		currentPostVote.VoteValue = utils.ConvertGQLVoteValueEnums(voteValue)
		_, err = currentPostVote.Update(ctx, database.DB, boil.Infer())
		if err != nil {
			return nil, err
		}
	}

	// return vote object
	gql_postVote := utils.ConvertPostVote(currentPostVote)
	return &gql_postVote, nil
}

func (r *mutationResolver) VoteOnComment(ctx context.Context, commentID int, voteValue model.VoteValue) (*model.CommentVote, error) {
	// authenticate user
	userID, err := middleware.GetUserIDFromSessions(ctx)
	if err != nil {
		return nil, err
	}

	// check if user has already voted
	currentCommentVote, err := sql_models.FindCommentVote(ctx, database.DB, commentID, userID)

	if err != nil {
		return nil, err
	}

	// attempt to add new vote or update existing vote
	if currentCommentVote == nil {
		currentCommentVote = &sql_models.CommentVote{
			CommentID: commentID,
			UserID:    userID,
			VoteValue: utils.ConvertGQLVoteValueEnums(voteValue),
		}
		err = currentCommentVote.Insert(ctx, database.DB, boil.Infer())
		if err != nil {
			return nil, err
		}

	} else {
		currentCommentVote.VoteValue = utils.ConvertGQLVoteValueEnums(voteValue)
		_, err = currentCommentVote.Update(ctx, database.DB, boil.Infer())
		if err != nil {
			return nil, err
		}
	}

	// return vote object
	gql_commentVote := utils.ConvertCommentVote(currentCommentVote)
	return &gql_commentVote, nil
}

/* -------------------------------------------------------------------------- */
/*                                manage users                                */
/* -------------------------------------------------------------------------- */

func (r *mutationResolver) RegisterNewUser(ctx context.Context, userInput model.UserInput) (*model.User, error) {
	// validate user inputs
	err := utils.ValidateEmail(userInput.Email)
	if err != nil {
		return nil, err
	}
	err = utils.ValidateUsername(userInput.Username)
	if err != nil {
		return nil, err
	}
	err = utils.ValidatePassword(userInput.Password)
	if err != nil {
		return nil, err
	}

	// get session
	_, session, err := middleware.GetGinContextAndSessions(ctx)
	if err != nil {
		return nil, err
	}

	// hash password
	hashedPassword, err := utils.HashPassword(userInput.Password)
	if err != nil {
		return nil, err
	}

	// insert new user
	newUser := sql_models.User{
		Username:     userInput.Username,
		Email:        userInput.Email,
		UserPassword: hashedPassword,
	}

	err = newUser.Insert(ctx, database.DB, boil.Infer())

	if err != nil {
		return nil, err
	}

	// format user and remove password from struct
	formattedUser := utils.ConvertUser(&newUser)

	// add new user to session
	session.Set("user", formattedUser.UserID)
	err = session.Save()
	if err != nil {
		return nil, err
	}

	return &formattedUser, err
}

func (r *mutationResolver) ToggleUserActiveStatus(ctx context.Context) (*model.User, error) {
	// authenticate user
	userID, err := middleware.GetUserIDFromSessions(ctx)
	if err != nil {
		return nil, err
	}

	// find current user status
	user, err := sql_models.FindUser(ctx, database.DB, userID)
	if err != nil {
		return nil, err
	}
	// toggle active status for user
	user.Active = !user.Active
	_, err = user.Update(ctx, database.DB, boil.Infer())
	if err != nil {
		return nil, err
	}

	// format as graphql response
	fmtUser := utils.ConvertUser(user)
	return &fmtUser, nil
}

func (r *mutationResolver) Login(ctx context.Context, username string, password string) (*model.User, error) {
	// get session
	_, session, err := middleware.GetGinContextAndSessions(ctx)
	if err != nil {
		return nil, err
	}
	// this function will accept either the username or user email
	user, err := sql_models.Users(qm.Where("(username = ?) or (email = ?)", username, username)).One(ctx, database.DB)
	if err != nil {
		return nil, err
	}

	// compare password with hashed password
	correctPassword := utils.CheckPasswordHash(password, user.UserPassword)
	if !correctPassword {
		return nil, errors.New(constants.INVALID_USERNAME_PASSWORD_ERROR_MESSAGE)
	}

	// access and save session
	session.Set("user", user.UserID)
	err = session.Save()
	if err != nil {
		return nil, err
	}

	// format user object and return it
	formattedUser := utils.ConvertUser(user)
	return &formattedUser, nil
}

func (r *mutationResolver) Logout(ctx context.Context) (bool, error) {
	// get gin context
	_, session, err := middleware.GetGinContextAndSessions(ctx)
	if err != nil {
		return false, err
	}

	// access and remove user from session
	session.Delete("user")
	err = session.Save()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) ForgotPassword(ctx context.Context, username string) (bool, error) {
	// confirm username / email correspond to user in DB
	user, err := sql_models.Users(qm.Where("username = ?", username), qm.Or("email = ?", username)).One(ctx, database.DB)
	if err != nil || user == nil {
		return false, err
	}

	// generate redis key using uuid that contains userid
	// which will be obtained through the url
	// and the url link shared via email with the user
	// so that only someone with access to the user email on record
	// should receive the reset link

	// generate unique redis key
	resetKeyUUID, err := uuid.NewV4()
	if err != nil {
		return false, err
	}

	resetKey := resetKeyUUID.String()
	// store user_id in redis using the unique key
	// with a one hour expiration
	_, err = database.RedisClient.Set(ctx, resetKey, user.UserID, time.Hour*1).Result()

	if err != nil {
		return false, err
	}

	// send email with custom password reset link to user's email
	err = utils.SendPasswordResetEmail(user.Email, resetKey)
	if err != nil {
		return false, err
	}

	// return true if successful
	return true, nil
}

func (r *mutationResolver) AccessPasswordReset(ctx context.Context, resetKey string) (bool, error) {
	// confirm the uuid-generated password-reset url is in our redis DB
	// before presenting reset form to the user
	_, err := database.RedisClient.Get(ctx, resetKey).Result()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) ResetPassword(ctx context.Context, resetKey string, userID int, newPassword string) (*model.User, error) {
	// confirm reset key is active in redis
	user_id, err := database.RedisClient.Get(ctx, resetKey).Result()
	if err != nil {
		return nil, err
	}

	err = utils.ValidatePassword(newPassword)
	if err != nil {
		return nil, err
	}

	// convert user_id to int
	user_id_int, err := strconv.Atoi(user_id)
	if err != nil {
		return nil, err
	}

	// locate user in database
	user, err := sql_models.FindUser(ctx, database.DB, user_id_int)
	if err != nil || user == nil {
		return nil, err
	}

	// hash new password
	hashedPassword, err := utils.HashPassword(newPassword)

	if err != nil {
		return nil, err
	}

	// update user password
	user.UserPassword = hashedPassword
	_, err = user.Update(ctx, database.DB, boil.Infer())
	if err != nil {
		return nil, err
	}

	// get gin context/ sessions and sign in user
	_, session, err := middleware.GetGinContextAndSessions(ctx)
	if err != nil {
		return nil, err
	}
	session.Set("user", user_id_int)
	err = session.Save()
	if err != nil {
		return nil, err
	}

	// convert sql user object to graphQL user object
	fmtUser := utils.ConvertUser(user)

	// return graphql user object
	return &fmtUser, nil
}

/* -------------------------------------------------------------------------- */
/*                       more dataloader field resolvers                      */
/* -------------------------------------------------------------------------- */

func (r *postResolver) User(ctx context.Context, obj *model.Post) (*model.User, error) {
	user, err := dataloader.For(ctx).UserById.Load(obj.UserID)
	return &user, err
}

func (r *postResolver) URLEncodedTitle(ctx context.Context, obj *model.Post) (string, error) {
	encodedTitle := url.QueryEscape(obj.Title)
	return encodedTitle, nil
}

func (r *postResolver) Comments(ctx context.Context, obj *model.Post) (*model.PaginatedComments, error) {
	paginatedComments, err := dataloader.For(ctx).CommentByPostID.Load(obj.PostID)
	return &paginatedComments, err
}

func (r *postResolver) Votes(ctx context.Context, obj *model.Post) (*model.Votes, error) {
	votes, err := dataloader.For(ctx).VotesByPostID.Load(obj.PostID)
	return &votes, err
}

/* -------------------------------------------------------------------------- */
/*                                  get posts                                 */
/* -------------------------------------------------------------------------- */

// get single post
func (r *queryResolver) GetPost(ctx context.Context, postID int) (*model.Post, error) {
	post, err := sql_models.Posts(qm.Where("post_id = ?", postID)).One(ctx, database.DB)

	if post == nil {
		return nil, err
	}

	// check if post is published
	// if unpublished, confirm user is author and authenticated
	if !post.Published {
		err = middleware.RejectIfNotAuthor(ctx, post.UserID)
		if err != nil {
			return nil, err
		}
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

func (r *queryResolver) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	user, err := sql_models.Users(qm.Where("username = ?", username)).One(ctx, database.DB)

	if user == nil {
		return nil, err
	}

	formattedUser := utils.ConvertUser(user)
	return &formattedUser, nil
}

// these will always be published since they will correspond to public urls
func (r *queryResolver) GetPostByUsernameAndTitle(ctx context.Context, username string, title string) (*model.Post, error) {
	// the title variable will pull a url-encoded slug, which will need to be unencoded before searching the database
	unencodedTitle, err := url.QueryUnescape(title)
	if err != nil {
		return nil, err
	}

	user, err := sql_models.Users(qm.Where("username = ? AND active = true", username)).One(ctx, database.DB)
	if err != nil {
		return nil, err
	}

	post, err := user.Posts(qm.Where("title = ? AND published = true", unencodedTitle)).One(ctx, database.DB)
	if err != nil {
		return nil, err
	}

	fmtPost := utils.ConvertPost(post)
	return &fmtPost, err
}

func (r *queryResolver) GetManyPosts(ctx context.Context, postSearch model.PostSearch, authorID int) (*model.PaginatedPosts, error) {
	// cap the maximum possible limit and return with one extra
	// to check for remaining posts
	var limitPlusOne int
	trueLimit := 20
	if postSearch.Limit > trueLimit {
		limitPlusOne = trueLimit + 1
	} else {
		limitPlusOne = postSearch.Limit + 1
	}

	// get posts from DB with optional search by title
	var posts sql_models.PostSlice
	if postSearch.Title == nil {
		retrievedPosts, err := sql_models.Posts(qm.Where("user_id = ? AND published = true", authorID), qm.Limit(limitPlusOne), qm.Offset(postSearch.Offset)).All(ctx, database.DB)
		if err != nil {
			return nil, err
		}
		posts = retrievedPosts
	} else {
		retrievedPosts, err := sql_models.Posts(qm.Where("user_id = ? AND published = true", authorID), qm.Limit(limitPlusOne), qm.Offset(postSearch.Offset), qm.Where("Title ILIKE ?", "%"+*postSearch.Title+"%")).All(ctx, database.DB)
		if err != nil {
			return nil, err
		}
		posts = retrievedPosts
	}

	// format posts for graphQL response
	formattedPosts := make([]*model.Post, len(posts))
	for i, value := range posts {
		fmtPost := utils.ConvertPost(value)
		formattedPosts[i] = &fmtPost
	}

	paginatedResponse := model.PaginatedPosts{
		Posts: formattedPosts,
		More:  len(posts) == limitPlusOne,
	}

	return &paginatedResponse, nil
}

// get unpublished posts for the current user - not visible to others
func (r *queryResolver) GetUnpublishedPosts(ctx context.Context, limit int, offset int) (*model.PaginatedPosts, error) {
	// get userID from sessions
	userID, err := middleware.GetUserIDFromSessions(ctx)
	if err != nil {
		return nil, err
	}

	// reject if unauthenticated
	if userID == 0 {
		return nil, errors.New(constants.ONLY_AUTHOR_ALLOWED_ERROR_MESSAGE)
	}

	// cap the maximum possible limit and return with one extra
	// to check for remaining posts
	var limitPlusOne int
	trueLimit := 20
	if limit > trueLimit {
		limitPlusOne = trueLimit + 1
	} else {
		limitPlusOne = limit + 1
	}

	// get unpublished posts for user
	// ignore title parameter of postInput, as these should be a fairly small number
	posts, err := sql_models.Posts(qm.Where("user_id = ? AND published = false", userID), qm.Limit(limitPlusOne), qm.Offset(offset)).All(ctx, database.DB)
	if err != nil {
		return nil, err
	}

	// format posts for graphQL response
	formattedPosts := make([]*model.Post, len(posts))
	for i, value := range posts {
		fmtPost := utils.ConvertPost(value)
		formattedPosts[i] = &fmtPost
	}

	paginatedResponse := model.PaginatedPosts{
		Posts: formattedPosts,
		More:  len(posts) == limitPlusOne,
	}

	return &paginatedResponse, nil
}

func (r *queryResolver) GetManyUsers(ctx context.Context, userSearch model.UserSearch) (*model.PaginatedUsers, error) {
	// cap the maximum possible limit and return with one extra
	// to check for remaining users
	var limitPlusOne int
	trueLimit := 20
	if userSearch.Limit > trueLimit {
		limitPlusOne = trueLimit + 1
	} else {
		limitPlusOne = userSearch.Limit + 1
	}

	// get users from DB with optional search by username
	var users sql_models.UserSlice
	if userSearch.Username == nil {
		retrievedUsers, err := sql_models.Users(qm.Limit(limitPlusOne), qm.Offset(userSearch.Offset)).All(ctx, database.DB)
		if err != nil {
			return nil, err
		}
		users = retrievedUsers
	} else {
		retrievedUsers, err := sql_models.Users(qm.Limit(limitPlusOne), qm.Offset(userSearch.Offset), qm.Where("username ILIKE %?%", userSearch.Username)).All(ctx, database.DB)
		if err != nil {
			return nil, err
		}
		users = retrievedUsers
	}

	// format users for graphQL response
	var formattedUsers []*model.User
	for _, value := range users {
		fmtUser := utils.ConvertUser(value)
		formattedUsers = append(formattedUsers, &fmtUser)
	}

	paginatedResponse := model.PaginatedUsers{
		Users: formattedUsers,
		More:  len(users) == limitPlusOne,
	}

	return &paginatedResponse, nil
}

/* -------------------------------------------------------------------------- */
/*                              get many comments                             */
/* -------------------------------------------------------------------------- */

func (r *queryResolver) GetManyComments(ctx context.Context, commentSearch model.CommentSearch) (*model.PaginatedComments, error) {
	// cap the maximum possible limit and return with one extra
	// to check for remaining comments
	var limitPlusOne int
	trueLimit := 20
	if commentSearch.Limit > trueLimit {
		limitPlusOne = trueLimit + 1
	} else {
		limitPlusOne = commentSearch.Limit + 1
	}

	// get comments from DB
	var whereClause string
	if commentSearch.ParentType == model.ParentTypePost {
		whereClause = "post_id = ?"
	} else {
		whereClause = "response_to_comment_id = ?"
	}
	retrievedComments, err := sql_models.Comments(qm.Where(whereClause, commentSearch.ParentID), qm.Limit(commentSearch.Limit), qm.Offset(commentSearch.Offset)).All(ctx, database.DB)
	if err != nil {
		return nil, err
	}

	// find which comments have subcomments
	var currentCommentIDList []int
	for _, value := range retrievedComments {
		currentCommentIDList = append(currentCommentIDList, value.CommentID)
	}

	// format ids as string for SQL ANY() argument
	queryParam := utils.FormatSliceForSQLParams(currentCommentIDList)
	subComments, err := sql_models.Comments(qm.Where("response_to_comment_id = ANY(?::int[])", queryParam)).All(ctx, database.DB)
	if err != nil {
		return nil, err
	}

	// format comments for graphQL response
	var formattedComments []*model.Comment
	for _, comment := range retrievedComments {
		// check if comment has subcomments
		hasSubComment := false
		for _, subComment := range subComments {
			if subComment.ResponseToCommentID.Int == comment.CommentID {
				hasSubComment = true
			}
		}
		fmtComment := utils.ConvertComment(comment, hasSubComment)
		formattedComments = append(formattedComments, &fmtComment)
	}

	paginatedResponse := model.PaginatedComments{
		Comments: formattedComments,
		More:     len(retrievedComments) == limitPlusOne,
	}

	return &paginatedResponse, nil
}

/* -------------------------------------------------------------------------- */
/*                          various utility functions    ??                     */
/* -------------------------------------------------------------------------- */

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	_, session, err := middleware.GetGinContextAndSessions(ctx)
	if err != nil {
		return nil, err
	}

	user := session.Get("user")
	if user == 0 || user == nil {
		return nil, nil
	}

	u, err := sql_models.Users(qm.Where("user_id = ?", user)).One(ctx, database.DB)
	if err != nil {
		//error log
		return nil, err
	}
	// format user and remove password from struct
	formattedUser := utils.ConvertUser(u)
	return &formattedUser, nil
}

func (r *queryResolver) IsAuthor(ctx context.Context, authorID int) (bool, error) {
	userID, err := middleware.GetUserIDFromSessions(ctx)
	if err != nil {
		return false, err
	}
	isAuthor := authorID == userID
	return isAuthor, nil
}

func (r *userResolver) Posts(ctx context.Context, obj *model.User) (*model.PaginatedPosts, error) {
	// since only the blog author is currently able to create posts
	// this should only be called for one user
	// and the dataloader pattern is currently not necessary

	// pagination will limit these to 20 posts
	// for fetching additional posts, the GetManyPosts resolver can then be used
	// with the limit and offset set accordingly
	posts, err := sql_models.Posts(qm.Where("user_id = ? AND published = true", obj.UserID), qm.Limit(21)).All(ctx, database.DB)
	if err != nil {
		return nil, err
	}

	// format posts for graphQL response
	var formattedPosts []*model.Post
	for _, value := range posts {
		fmtPost := utils.ConvertPost(value)
		formattedPosts = append(formattedPosts, &fmtPost)
	}

	// check if there are more posts and remove the final one
	// if the full 21 posts were retrieved
	hasMore := false
	if len(formattedPosts) == 21 {
		hasMore = true
		formattedPosts = formattedPosts[:len(formattedPosts)-1]
	}

	response := model.PaginatedPosts{
		Posts: formattedPosts,
		More:  hasMore,
	}

	return &response, nil
}

func (r *userResolver) Comments(ctx context.Context, obj *model.User) (*model.PaginatedComments, error) {
	paginatedComments, err := dataloader.For(ctx).CommentByUserID.Load(obj.UserID)
	return &paginatedComments, err
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
