package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/antsrp/banner_service/internal/domain/models"
	"github.com/antsrp/banner_service/internal/repository"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserStorage struct {
	conn *Connection
}

func NewUserStorage(conn *Connection) UserStorage {
	return UserStorage{
		conn: conn,
	}
}

func (s UserStorage) Create(ctx context.Context, user models.User) repository.DatabaseError {
	tx, err := s.conn.PC.Begin(ctx)
	if err != nil {
		return NewError("can't create transaction", err)
	}
	defer tx.Rollback(ctx)
	var id int
	if err := tx.QueryRow(ctx, `INSERT INTO users (name, is_admin) VALUES ($1, $2) RETURNING id;`, user.Name, user.IsAdmin).
		Scan(&id); err != nil {
		errString := "can't create user"
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			err = repository.ErrUsernameAlreadyExists
		}
		return NewError(errString, err)
	}
	values := make([]string, 0, len(user.Tags))
	for _, tag := range user.Tags {
		values = append(values, fmt.Sprintf("(%d, %d)", id, tag))
	}
	if _, err := tx.Exec(ctx, fmt.Sprintf(`INSERT INTO users_tags (user_id, tag_id) VALUES %s`, strings.Join(values, ","))); err != nil {
		return NewError("can't add tags for user", err)
	}

	tx.Commit(ctx)
	return nil
}
func (s UserStorage) FindByName(ctx context.Context, name string) (repository.UserWithToken, repository.DatabaseError) {
	var uwt repository.UserWithToken
	var token sql.NullString
	query := `SELECT users.id, name, is_admin, token FROM users 
	LEFT JOIN tokens ON users.id = tokens.user_id 
	WHERE name = $1`
	if err := s.conn.PC.QueryRow(ctx, query, name).Scan(&uwt.ID, &uwt.Name, &uwt.IsAdmin, &token); err != nil {
		errString := "can't find user by name"
		if errors.Is(err, pgx.ErrNoRows) {
			err = repository.ErrEntityNotFound
		}
		return repository.UserWithToken{}, NewError(errString, err)
	}
	if token.Valid {
		uwt.Token = token.String
	}
	return uwt, nil
}
func (s UserStorage) AddToken(ctx context.Context, uwt repository.UserWithToken) repository.DatabaseError {
	if _, err := s.conn.PC.Exec(ctx, `INSERT INTO tokens (user_id, token) VALUES ($1, $2)
	ON CONFLICT (user_id)
	DO UPDATE SET token = $3, created_at = now()`, uwt.ID, uwt.Token, uwt.Token); err != nil {
		return NewError("can't add token to user", err)
	}
	return nil
}

var _ repository.UserStorage = UserStorage{}
