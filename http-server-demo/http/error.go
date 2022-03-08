package http

import (
	"fmt"
	"io"
	"net/http"
)

type Error struct {
	code    int
	message string
}

// NewError returns a new Error with the given
// HTTP status code and error message.
//
// Two errors with the same status code and
// error message are equal.
func NewError(code int, msg string) Error {
	return Error{
		code:    code,
		message: msg,
	}
}

// Status returns the HTTP status code of the error.
func (e Error) Status() int { return e.code }

func (e Error) Error() string { return e.message }

// Error sends the given err as JSON error response to w.
//
// If err has a 'Status() int' method then Error sets the
// response status code to err.Status(). Otherwise, it will
// send 500 (internal server error).
//
// If err is nil then Error will send the status code 500 and
// an empty JSON response body - i.e. '{}'.
func sendHttpResError(w http.ResponseWriter, err error) error {
	var status = http.StatusInternalServerError
	if e, ok := err.(interface{ Status() int }); ok {
		status = e.Status()
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)

	const (
		emptyMsg = `{}`
		format   = `{"message":"%v"}`
	)
	if err == nil {
		_, err = io.WriteString(w, emptyMsg)
	} else {
		_, err = io.WriteString(w, fmt.Sprintf(format, err))
	}
	return err
}
