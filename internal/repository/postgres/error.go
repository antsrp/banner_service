package postgres

import (
	"errors"
	"fmt"

	"github.com/antsrp/banner_service/internal/repository"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type DatabaseError struct {
	isInternal bool
	real       error
}

func NewError(realErrorString string, realError error) DatabaseError {
	var isInternal bool
	var pgErr *pgconn.PgError
	if errors.As(realError, &pgErr) {
		isInternal = pgerrcode.IsInternalError(pgErr.Code)
	}
	return DatabaseError{
		isInternal: isInternal,
		real:       fmt.Errorf("%s: %w", realErrorString, realError),
	}
}

func (d DatabaseError) IsInternal() bool {
	return d.isInternal
}

func (d DatabaseError) Cause() error {
	return d.real
}

var _ repository.DatabaseError = DatabaseError{}
