package constants

// list different types of custom error messages
var UNAUTHENTICATED_ERROR_MESSAGE = "Must be logged in!"
var ONLY_AUTHOR_ALLOWED_ERROR_MESSAGE = "Only the author of the blog may add, edit, or delete posts"
var ONLY_COMMENT_AUTHOR_MAY_EDIT = "Only the author of a comment can edit, delete, or restore it"
var PASSWORD_TOO_SHORT_ERROR_MESSAGE = "Password must be 8 or more characters long"
var PASSWORD_LACKS_MIX_OF_CHARS_ERROR_MESSAGE = "Password must have letters, numbers, and special characters"
var PASSWORD_LACKS_UPPER_AND_LOWERCASE_LETTERS_ERROR_MESSAGE = "Password must have both lower and uppercase characters"
var USERNAME_TOO_SHORT_ERROR_MESSAGE = "Username must be at least 3 characters long"
var USERNAME_INAPPROPRIATE_ERROR_MESSAGE = "Must use an appropriate username"
var USERNAME_NOT_URL_COMPATIBLE_ERROR_MESSAGE = "Please use only letters, numbers, and hyphens for your username"
var INVALID_EMAIL_ERROR_MESSAGE = "Must use a valid email address"
var INVALID_USERNAME_PASSWORD_ERROR_MESSAGE = "Incorrect username / password combination!"

// confirm if error has custom error message
// which can be shared directly with the client
func IsCustomError(errMessage string) bool {
	// store custom errors in array
	errorMessages := [...]string{
		UNAUTHENTICATED_ERROR_MESSAGE,
		ONLY_AUTHOR_ALLOWED_ERROR_MESSAGE,
		ONLY_COMMENT_AUTHOR_MAY_EDIT,
		PASSWORD_TOO_SHORT_ERROR_MESSAGE,
		PASSWORD_LACKS_MIX_OF_CHARS_ERROR_MESSAGE,
		PASSWORD_LACKS_UPPER_AND_LOWERCASE_LETTERS_ERROR_MESSAGE,
		USERNAME_TOO_SHORT_ERROR_MESSAGE,
		USERNAME_INAPPROPRIATE_ERROR_MESSAGE,
		USERNAME_NOT_URL_COMPATIBLE_ERROR_MESSAGE,
		INVALID_EMAIL_ERROR_MESSAGE,
		INVALID_USERNAME_PASSWORD_ERROR_MESSAGE,
	}

	// loop through to find match
	for _, value := range errorMessages {
		if errMessage == value {
			return true
		}
	}

	// return false if no match found
	return false
}