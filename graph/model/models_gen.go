// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type Comment struct {
	CommentID           int                `json:"comment_id"`
	ResponseToCommentID *int               `json:"response_to_comment_id"`
	PostID              int                `json:"post_id"`
	UserID              int                `json:"user_id"`
	User                *User              `json:"user"`
	CommentText         string             `json:"comment_text"`
	CreatedAt           time.Time          `json:"created_at"`
	Comments            *PaginatedComments `json:"comments"`
	Votes               *Votes             `json:"votes"`
	Deleted             bool               `json:"deleted"`
	HasSubComments      bool               `json:"hasSubComments"`
}

type CommentSearch struct {
	ParentID   int        `json:"parent_id"`
	ParentType ParentType `json:"parent_type"`
	Offset     int        `json:"offset"`
	Limit      int        `json:"limit"`
}

type CommentVote struct {
	CommentID int       `json:"comment_id"`
	VoteValue VoteValue `json:"vote_value"`
	UserID    int       `json:"user_id"`
}

type PaginatedComments struct {
	Comments []*Comment `json:"comments"`
	More     bool       `json:"more"`
}

type PaginatedPosts struct {
	Posts []*Post `json:"posts"`
	More  bool    `json:"more"`
}

type PaginatedUsers struct {
	Users []*User `json:"users"`
	More  bool    `json:"more"`
}

type Post struct {
	PostID          int                `json:"post_id"`
	UserID          int                `json:"user_id"`
	User            *User              `json:"user"`
	Title           string             `json:"title"`
	URLEncodedTitle string             `json:"urlEncodedTitle"`
	Subtitle        string             `json:"subtitle"`
	PostText        string             `json:"post_text"`
	CreatedAt       time.Time          `json:"created_at"`
	Comments        *PaginatedComments `json:"comments"`
	Votes           *Votes             `json:"votes"`
	Deleted         bool               `json:"deleted"`
	Published       bool               `json:"published"`
}

type PostInput struct {
	Title     string  `json:"title"`
	Subtitle  *string `json:"subtitle"`
	Text      string  `json:"text"`
	Published bool    `json:"published"`
}

type PostSearch struct {
	Title  *string `json:"title"`
	Limit  int     `json:"limit"`
	Offset int     `json:"offset"`
}

type PostVote struct {
	PostID    int       `json:"post_id"`
	VoteValue VoteValue `json:"vote_value"`
	UserID    int       `json:"user_id"`
}

type User struct {
	UserID    int                `json:"user_id"`
	Username  string             `json:"username"`
	Email     string             `json:"email"`
	Posts     *PaginatedPosts    `json:"posts"`
	Comments  *PaginatedComments `json:"comments"`
	CreatedAt time.Time          `json:"created_at"`
	Active    bool               `json:"active"`
}

type UserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserSearch struct {
	Username *string `json:"username"`
	Limit    int     `json:"limit"`
	Offset   int     `json:"offset"`
}

type Votes struct {
	Upvote   int `json:"upvote"`
	Downvote int `json:"downvote"`
}

type ParentType string

const (
	ParentTypePost    ParentType = "post"
	ParentTypeComment ParentType = "comment"
)

var AllParentType = []ParentType{
	ParentTypePost,
	ParentTypeComment,
}

func (e ParentType) IsValid() bool {
	switch e {
	case ParentTypePost, ParentTypeComment:
		return true
	}
	return false
}

func (e ParentType) String() string {
	return string(e)
}

func (e *ParentType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ParentType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ParentType", str)
	}
	return nil
}

func (e ParentType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type VoteValue string

const (
	VoteValueUpvote   VoteValue = "upvote"
	VoteValueDownvote VoteValue = "downvote"
	VoteValueNeutral  VoteValue = "neutral"
)

var AllVoteValue = []VoteValue{
	VoteValueUpvote,
	VoteValueDownvote,
	VoteValueNeutral,
}

func (e VoteValue) IsValid() bool {
	switch e {
	case VoteValueUpvote, VoteValueDownvote, VoteValueNeutral:
		return true
	}
	return false
}

func (e VoteValue) String() string {
	return string(e)
}

func (e *VoteValue) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = VoteValue(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid VoteValue", str)
	}
	return nil
}

func (e VoteValue) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
