package example

import "net/http"

// Inspiration from:
// https://www.joeshaw.org/error-handling-in-go-http-applications/

// APIError ...
type APIError interface {
	// APIError returns an HTTP status code and an API-safe error message.
	APIError() (int, string)
}

type sentinelAPIError struct {
	status int
	msg    string
}

func (e sentinelAPIError) Error() string {
	return e.msg
}

func (e sentinelAPIError) APIError() (int, string) {
	return e.status, e.msg
}

var (
	ErrAuth       = &sentinelAPIError{status: http.StatusUnauthorized, msg: "invalid token"}
	ErrValidation = &sentinelAPIError{status: http.StatusBadRequest, msg: "invalid request"}
	ErrNotFound   = &sentinelAPIError{status: http.StatusNotFound, msg: "not found"}
	ErrTemporary  = &sentinelAPIError{status: http.StatusServiceUnavailable, msg: "temporary error"}
)

// sentinelWrappedError....
type sentinelWrappedError struct {
	error
	sentinel *sentinelAPIError
}

func (e sentinelWrappedError) Is(err error) bool {
	return e.sentinel == err
}

func (e sentinelWrappedError) APIError() (int, string) {
	return e.sentinel.APIError()
}

func WrapError(err error, sentinel *sentinelAPIError) error {
	return sentinelWrappedError{error: err, sentinel: sentinel}
}
