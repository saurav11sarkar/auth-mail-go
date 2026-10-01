package user

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrUpdateNotImplemented = errors.New("profile update is not implemented")
)
