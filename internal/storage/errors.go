package storage

import "errors"

var (
	ErrConflict       = errors.New("entity conflict")
	ErrBatchIsEmpty   = errors.New("batch is empty")
	ErrURLNotValid    = errors.New("url is invalidate")
	ErrNotFound       = errors.New("entity not found")
	ErrDeleteAccepted = errors.New("entity accepted delete")
)
