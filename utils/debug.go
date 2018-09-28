package utils

// Error is a type for declaring a response error from server
type Error struct {
	Code int // status code of the error
	Err  error
}

// Reason is a method for type `Error` which returns the reson for the non 200 status code
func (e *Error) Reason() string {
	return e.Err.Error()
}
