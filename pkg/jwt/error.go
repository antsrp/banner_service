package jwt

import "errors"

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrUnexpectedMethod = errors.New("unexpected method")
)
