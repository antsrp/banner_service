package repository

import "database/sql"

type Connection interface {
	Check() error
	Close() error
	DB() (*sql.DB, error)
}
