package repository

import "errors"

const (
	msgNoRowsAffected = "no rows affected"
	msgEntityNotFound = "no entity found"

	msgUsernameAlreadyExists = "user with name already exists"
)

var (
	ErrNoRowsAffected = errors.New(msgNoRowsAffected)
	ErrEntityNotFound = errors.New(msgEntityNotFound)

	ErrUsernameAlreadyExists = errors.New(msgUsernameAlreadyExists)
)

type DatabaseError interface {
	IsInternal() bool
	Cause() error
}
