package domain

import "errors"

var (
	ErrLinkNotFound           = errors.New("short code not found")
	ErrCodeUniquenessConflict = errors.New("provided short code already exists")
	ErrInternal               = errors.New("something went wrong")
	ErrNoURLProvided          = errors.New("url required")
	ErrNoCodeProvided         = errors.New("code required")
)
