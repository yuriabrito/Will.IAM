package errors

// ErrorWithStatusCode is used to determine whether an error
// has a defined status code for an http response
type ErrorWithStatusCode interface {
	Error() string
	StatusCode() int
}
