CREATE DATABASE clean_blog; 

CREATE TABLE users (
  userID SERIAL PRIMARY KEY,
  username VARCHAR(255) UNIQUE NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  user_password VARCHAR(255) NOT NULL,
  date_joined TIMESTAMPTZ NOT NULL
);

CREATE TABLE posts (
  postID SERIAL PRIMARY KEY,
  userID INT REFERENCES Users(userID) NOT NULL,
  title VARCHAR(255) NOT NULL,
  subTitle VARCHAR(255), -- nullable
  post_text TEXT NOT NULL, -- may change to JSONB based on react editor
  date_posted TIMESTAMPTZ NOT NULL
);

CREATE TABLE comments (
  commentID SERIAL PRIMARY KEY,
  responseToCommentID INT REFERENCES Comments (commentID), -- nullable, used when one comment is in response to another comment, nesting it
  postID INT REFERENCES Posts(postID) NOT NULL,
  userID INT REFERENCES Users(userID) NOT NULL,
  comment_text TEXT NOT NULL,
  date_posted TIMESTAMPTZ NOT NULL
);

CREATE TABLE post_votes (
  postID INT REFERENCES Posts(postID) NOT NULL,
  voteValue INT NOT NULL, -- 1, 0, or -1
  userID INT REFERENCES Users(userID) NOT NULL
);

CREATE TABLE comment_votes (
  commentID INT REFERENCES Comments(commentID) NOT NULL,
  voteValue INT NOT NULL, -- 1, 0, or -1
  userID INT REFERENCES Users(userID) NOT NULL
);