package utils

import (
	"errors"
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrTokenNotFound       = errors.New("token not found")
	ErrTokenExpired        = errors.New("token expired")
	ErrTokenRevoked        = errors.New("token revoked")
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrInvalidToken        = errors.New("invalid token")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrUserNotVerified     = errors.New("user account not verified by admin")
	ErrUserAlreadyVerified = errors.New("user already verified")
	ErrUnauthorizedAdmin   = errors.New("admin access required")
)

func IsValidationError(err error) bool {
	return err != nil && (err.Error() == "validation failed" ||
		err == ErrInvalidCredentials ||
		err == ErrUserAlreadyExists)
}

func IsNotFoundError(err error) bool {
	return err == ErrUserNotFound || err == ErrTokenNotFound
}

func IsAuthError(err error) bool {
	return err == ErrInvalidCredentials ||
		err == ErrTokenExpired ||
		err == ErrTokenRevoked ||
		err == ErrInvalidToken ||
		err == ErrUnauthorized
}
