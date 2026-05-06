package apperrors

// Account Errors
const (
	DuplicateEmail     = "An account with that email already exists"
	InvalidCredentials = "Invalid email and password combination"
)

// Generic Errors
const (
	InvalidSession = "Provided session is invalid"
	ServerError    = "Something went wrong. Try again later"
)
