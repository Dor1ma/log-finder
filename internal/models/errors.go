package models

import "errors"

var (
	ErrNotFound      = errors.New("log entry not found")
	ErrInvalidFormat = errors.New("invalid log format")
)
