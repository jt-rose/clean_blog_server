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
var INVALID_EMAIL_ERROR_MESSAGE = "Must use a valid email address"

// store custom errors in array
// which will be used by the errorHandler
func createErrMsgArray() []string {
	errorMessages := [...]string{
		UNAUTHENTICATED_ERROR_MESSAGE,
		ONLY_AUTHOR_ALLOWED_ERROR_MESSAGE,
		ONLY_COMMENT_AUTHOR_MAY_EDIT,
		PASSWORD_TOO_SHORT_ERROR_MESSAGE,
		PASSWORD_LACKS_MIX_OF_CHARS_ERROR_MESSAGE,
		PASSWORD_LACKS_UPPER_AND_LOWERCASE_LETTERS_ERROR_MESSAGE,
		USERNAME_TOO_SHORT_ERROR_MESSAGE,
		USERNAME_INAPPROPRIATE_ERROR_MESSAGE,
		INVALID_EMAIL_ERROR_MESSAGE,
	}

	return errorMessages[:]
}

var ErrorMessages = createErrMsgArray()