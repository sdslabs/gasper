package types

import (
	"fmt"
	"strings"
)

// ResponseError is a type for declaring a response error from server
type ResponseError struct {
	Code int    // status code of the error
	Msg  string // Default set to Err.Error()
}

// NewResponseError returns a ResponseError type with the given message
// Alternatively, you can provinde an error interface that can replace the Msg
func NewResponseError(code int, msg string, err error) *ResponseError {
	if strings.Trim(msg, " ") != "" {
		return &ResponseError{
			Code: code,
			Msg:  msg,
		}
	} else if err != nil {
		return &ResponseError{
			Code: code,
			Msg:  err.Error(),
		}
	}
	return nil
}

// Error is the method to implement error interface of Golang
func (err *ResponseError) Error() string {
	return fmt.Sprintf("Status %d: %s", err.Code, err.Msg)
}

// Reason returns the message behind the error
func (err *ResponseError) Reason() string {
	return err.Msg
}

// Status returns server response code for the error
func (err *ResponseError) Status() int {
	return err.Code
}
