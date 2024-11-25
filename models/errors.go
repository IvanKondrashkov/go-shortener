package models

import (
	"errors"
)

var (
	ErrNotFound = errors.New("entity not found")
	ErrConflict = errors.New("entity id conflict")
)
