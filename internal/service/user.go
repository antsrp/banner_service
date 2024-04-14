package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/antsrp/banner_service/internal/domain/models"
	"github.com/antsrp/banner_service/internal/repository"
	"github.com/antsrp/banner_service/pkg/jwt"
	"github.com/antsrp/banner_service/pkg/logger"
)

type UserStorager interface {
	UserByToken(string) (models.User, Error)
	GenerateToken(models.User) (string, Error)
	SignIn(string) (string, Error)
}

type UserService struct {
	userStorage repository.UserStorage
	jwtService  jwt.Service
	logger      logger.Logger
}

func NewUserService(us repository.UserStorage, js jwt.Service, logger logger.Logger) UserService {
	return UserService{
		userStorage: us,
		jwtService:  js,
		logger:      logger,
	}
}

func (s UserService) UserByToken(token string) (models.User, Error) {
	data, err := s.jwtService.Parse(token)
	if err != nil {
		s.logger.Info(fmt.Errorf("error while parsing token: %w", err).Error())
		return models.User{}, NewServiceError(false, jwt.ErrInvalidToken)
	}
	var user models.User

	if name, ok := data[`username`].(string); ok {
		user.Name = name
	} else {
		return models.User{}, NewServiceError(false, fmt.Errorf("bad token"))
	}

	if isAdmin, ok := data[`is_admin`].(bool); ok {
		user.IsAdmin = isAdmin
	} else {
		return models.User{}, NewServiceError(false, fmt.Errorf("bad token"))
	}

	return user, nil
}

func (s UserService) GenerateToken(user models.User) (string, Error) {
	token, err := s.jwtService.NewToken(map[string]any{
		`is_admin`:   user.IsAdmin,
		`username`:   user.Name,
		`created_at`: time.Now().Unix(),
	})
	if err != nil {
		s.logger.Info(fmt.Errorf("can't create token: %w", err).Error())
		return "", defaultInternalError
	}

	return token, nil
}

func (s UserService) SignIn(name string) (string, Error) {
	ctx := context.Background()
	user, err := s.userStorage.FindByName(ctx, name)
	if err != nil {
		s.logger.Info(fmt.Errorf("can't find user: %w", err.Cause()).Error())
		if errors.Is(err.Cause(), repository.ErrEntityNotFound) {
			return "", NewServiceError(false, fmt.Errorf("user not found"))
		}
		return "", defaultInternalError
	}
	if user.Token == "" {
		token, err := s.jwtService.NewToken(map[string]any{
			`is_admin`:   user.IsAdmin,
			`username`:   user.Name,
			`created_at`: time.Now().Unix(),
		})
		if err != nil {
			s.logger.Error(fmt.Errorf("can't create token: %w", err).Error())
			return "", defaultInternalError
		}
		user.Token = token
	}
	if err := s.userStorage.AddToken(ctx, user); err != nil {
		s.logger.Error(fmt.Errorf("can't add token to storage: %w", err.Cause()).Error())
		if err.IsInternal() {
			return "", defaultInternalError
		}
		return "", NewServiceError(false, err.Cause())
	}

	return user.Token, nil
}

var _ UserStorager = UserService{}
