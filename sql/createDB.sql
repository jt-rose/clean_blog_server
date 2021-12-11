CREATE DATABASE clean_blog; 

CREATE TABLE users (
  user_id SERIAL PRIMARY KEY,
  username VARCHAR(255) UNIQUE NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  user_password VARCHAR(255) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE posts (
  post_id SERIAL PRIMARY KEY,
  user_id INT REFERENCES Users(user_id) NOT NULL,
  title VARCHAR(255) NOT NULL,
  subtitle VARCHAR(255), NOT NULL,
  post_text TEXT NOT NULL, -- may change to JSONB based on react editor
  created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE comments (
  comment_id SERIAL PRIMARY KEY,
  response_to_comment_id INT REFERENCES Comments (comment_id), -- nullable, used when one comment is in response to another comment, nesting it
  post_id INT REFERENCES Posts(post_id) NOT NULL,
  user_id INT REFERENCES Users(user_id) NOT NULL,
  comment_text TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE post_votes (
  post_id INT REFERENCES Posts(post_id) NOT NULL,
  vote_value INT NOT NULL, -- 1, 0, or -1
  user_id INT REFERENCES Users(user_id) NOT NULL,
  PRIMARY KEY(post_id, user_id)
);

CREATE TABLE comment_votes (
  comment_id INT REFERENCES Comments(comment_id) NOT NULL,
  vote_value INT NOT NULL, -- 1, 0, or -1
  user_id INT REFERENCES Users(user_id) NOT NULL,
  PRIMARY KEY(comment_id, user_id)
);

CREATE TABLE error_log (
    log_id SERIAL PRIMARY KEY,
    err_message TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    error_origin VARCHAR(255) NOT NULL
);