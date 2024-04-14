package repository

import (
	"context"

	"github.com/antsrp/banner_service/internal/domain/models"
)

type UserWithToken struct {
	ID int
	models.User
	Token string
}

type UserStorage interface {
	Create(context.Context, models.User) DatabaseError
	FindByName(context.Context, string) (UserWithToken, DatabaseError)
	AddToken(context.Context, UserWithToken) DatabaseError
}
