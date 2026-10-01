package user

import "errors"

var (
	ErrInvalidListQuery     = errors.New("invalid user list query")
	ErrUserNotFound         = errors.New("user not found")
	ErrUpdateNotImplemented = errors.New("profile update is not implemented")
)
