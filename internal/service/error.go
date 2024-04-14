package service

import "fmt"

type Error interface {
	IsInternal() bool
	Cause() error
}

type serviceError struct {
	isInternal bool
	real       error
}

var defaultInternalError = serviceError{isInternal: true, real: ErrDefaultInternalError}

func (e serviceError) IsInternal() bool {
	return e.isInternal
}

func (e serviceError) Cause() error {
	return e.real
}
func NewServiceError(isInternal bool, err error) Error {
	return serviceError{
		isInternal: isInternal,
		real:       err,
	}
}

var (
	ErrDefaultInternalError = fmt.Errorf("internal server error, try again later")
	ErrBannerNotFound       = fmt.Errorf("banner not found")
)
