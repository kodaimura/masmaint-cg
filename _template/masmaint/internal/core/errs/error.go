package errs

import (
	"fmt"
)


/////////////////////////////////////////////////////////////////////////
type BadRequestError struct {
	Field string
}

func NewBadRequestError(field string) error {
	return BadRequestError{Field: field}
}

func (e BadRequestError) Error() string {
	if e.Field == "" {
		return "error: The content of the request is invalid."
	}
	return fmt.Sprintf("error: Field '%s' binding failed.", e.Field)
}

/////////////////////////////////////////////////////////////////////////
type UniqueConstraintError struct {
	Column string
}

func NewUniqueConstraintError(column string) error {
	return UniqueConstraintError{Column: column}
}

func (e UniqueConstraintError) Error() string {
	return fmt.Sprintf("error: Column '%s' should be unique.", e.Column)
}

/////////////////////////////////////////////////////////////////////////
type NotFoundError struct {}

func NewNotFoundError() error {
	return NotFoundError{}
}

func (e NotFoundError) Error() string {
	return "error: Not found"
}

/////////////////////////////////////////////////////////////////////////
type UnexpectedError struct {
	Message string
}

func NewUnexpectedError(message string) error {
	return UnexpectedError{Message: message}
}

func (e UnexpectedError) Error() string {
	return fmt.Sprintf("error: %s", e.Message)
}