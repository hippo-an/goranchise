package services

type InvalidPasswordTokenError struct{}

func (e InvalidPasswordTokenError) Error() string {
	return "invalid token"
}

type NotAuthenticatedError struct {
}

func (e NotAuthenticatedError) Error() string {
	return "user not authenticated"
}
