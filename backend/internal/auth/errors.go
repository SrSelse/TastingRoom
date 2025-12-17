package auth

type UnauthenticatedError struct {
	ErrorInfo string
}

type UserNotFoundError struct {
	ErrorInfo string
}

type DatabaseError struct {
	Err error
}

func (e UnauthenticatedError) Error() string {
	return e.ErrorInfo
}

func (e UserNotFoundError) Error() string {
	return e.ErrorInfo
}

func (e DatabaseError) Error() string {
	return "DatabaseError"
}
