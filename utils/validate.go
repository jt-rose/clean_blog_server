package utils

import (
	"errors"
	"net/url"
	"regexp"

	"github.com/jt-rose/clean_blog_server/constants"

	goaway "github.com/TwiN/go-away"
)

func ValidateEmail(email string) error {
	// match *@*.*
	match, _ := regexp.MatchString(`^\w+@\w+\.\w+$`, email)
	if match {
		return errors.New(constants.INVALID_EMAIL_ERROR_MESSAGE)
	}
	return nil
}

func ValidateUsername(username string) error {
	if len(username) < 3 {
		return errors.New(constants.USERNAME_TOO_SHORT_ERROR_MESSAGE)
	}
	if goaway.IsProfane(username) {
		return errors.New(constants.USERNAME_INAPPROPRIATE_ERROR_MESSAGE)
	}
	// the username will be used in a url and should be compatible with a url query
	if url.QueryEscape(username) != username {
		return errors.New(constants.USERNAME_NOT_URL_COMPATIBLE_ERROR_MESSAGE)
	}
	return nil
}

func ValidatePassword(password string) error {
	// must be 8 characters or more
	if len(password) < 8 {
		return errors.New(constants.PASSWORD_TOO_SHORT_ERROR_MESSAGE)
	}
	// must include letter, number, and special character
	letter, _ := regexp.MatchString(`[a-zA-Z]`, password)
	num, _ := regexp.MatchString(`[0-9]`, password)
	specChar, _ := regexp.MatchString(`[!@#$%&*?]`, password)
	if !letter || !num || !specChar {
		return errors.New(constants.PASSWORD_LACKS_UPPER_AND_LOWERCASE_LETTERS_ERROR_MESSAGE)
	}
	// mix of upper and lowercase
	lower, _ := regexp.MatchString(`[a-z]`, password)
	upper, _ := regexp.MatchString(`[A-Z]`, password)
	if !lower || !upper {
		return errors.New(constants.PASSWORD_LACKS_UPPER_AND_LOWERCASE_LETTERS_ERROR_MESSAGE)
	}

	return nil
}

