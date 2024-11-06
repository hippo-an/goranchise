package container

type InvalidTokenError struct{}

func (e InvalidTokenError) Error() string {
	return "invalid token"
}

type NotAuthenticatedError struct {
}

func (e NotAuthenticatedError) Error() string {
	return "user not authenticated"
}
