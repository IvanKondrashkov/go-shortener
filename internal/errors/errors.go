package errors

import (
	"errors"
)

var (
	ErrConflict    = errors.New("entity conflict")
	ErrURLNotValid = errors.New("url is invalidate")
	ErrNotFound    = errors.New("entity not found")
)
