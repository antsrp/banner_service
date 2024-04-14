package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/antsrp/banner_service/internal/repository"
	"github.com/antsrp/banner_service/pkg/infrastructure/db"
	"github.com/antsrp/banner_service/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type connection struct {
	PC     *pgxpool.Pool
	logger logger.Logger
	ctx    context.Context
}

func NewConnection(ctx context.Context, settings db.Settings, logger logger.Logger) (*Connection, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		settings.User, settings.Password, settings.Host, settings.Port, settings.DBName)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("can't create connection pool: %w", err)
	}

	c := connection{
		PC:     pool,
		logger: logger,
		ctx:    ctx,
	}

	if err := c.Check(); err != nil {
		return nil, err
	}
	c.logger.Info("postgres connection opened")

	return &Connection{c}, nil
}

func (c connection) Check() error {
	return c.PC.Ping(c.ctx)
}

func (c connection) Close() error {
	c.logger.Info("postgres connection closing")
	c.PC.Close()
	return nil
}

func (c connection) DB() (*sql.DB, error) {
	return stdlib.OpenDBFromPool(c.PC), nil
}

type Connection struct {
	connection
}

var _ repository.Connection = Connection{}
