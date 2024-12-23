package server

import (
	"fmt"
)

// Error - custom server error.
type Error struct {
	err error
}

func (e *Error) Error() string {
	return fmt.Sprintf("[server]: %v", e.err)
}

func NewError(err error) error {
	if err == nil {
		return nil
	}

	return &Error{
		err: err,
	}
}

func (e *Error) Unwrap() error {
	return e.err
}
