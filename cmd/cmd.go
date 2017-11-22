package cmd

import (
	"fmt"
)

// ErrMissingParameters this is error the parameter is missing
type ErrMissingParameters struct {
	Param string
}

// Error return the message error
func (mp ErrMissingParameters) Error() string {
	return fmt.Sprintf("parameter %s is missing", mp.Param)
}

// ErrBodyEmpty return message that body cannot be empty
type ErrBodyEmpty struct{}

// Error return the message error
func (be ErrBodyEmpty) Error() string {
	return fmt.Sprintf("body cannot be empty")
}
