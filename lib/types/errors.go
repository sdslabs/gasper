package types

import (
	"fmt"
	"strings"
)

// ResponseError is a type for declaring response error from server
type ResponseError interface {
	error
	Message() string
	Verbose() string
	Status() int
}

// for implementing ResponseError
type rerror struct {
	Code int
	Msg  string
	Err  error
}

// NewResErr returns a ResponseError type with the given message
// Alternatively, you can provinde an error interface that can replace the Msg
// Both msg and error cannot be nil (empty)
func NewResErr(code int, msg string, err error) ResponseError {
	var message string
	if strings.Trim(msg, " ") != "" {
		message = msg
	} else { // Assuming if msg is empty, err is not nil
		message = err.Error()
	}
	return &rerror{
		Code: code,
		Msg:  message,
		Err:  err,
	}
}

// Error is the method to implement error interface of Golang
func (err *rerror) Error() string {
	return fmt.Sprintf("%d: %s\n%s", err.Code, err.Msg, err.Err.Error())
}

// Message returns the message accompanying the error
func (err *rerror) Message() string {
	return err.Msg
}

// Verbose returns the reason behind the error
func (err *rerror) Verbose() string {
	return err.Err.Error()
}

// Status returns server response code for the error
func (err *rerror) Status() int {
	return err.Code
}
