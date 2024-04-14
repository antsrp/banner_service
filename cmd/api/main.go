package main

import (
	"context"
	"os"

	"github.com/antsrp/banner_service/internal/cache/redis"
	"github.com/antsrp/banner_service/internal/domain/models"
	"github.com/antsrp/banner_service/internal/repository/postgres"
	"github.com/antsrp/banner_service/internal/rest"
	"github.com/antsrp/banner_service/internal/service"
	"github.com/antsrp/banner_service/pkg/config"
	cs "github.com/antsrp/banner_service/pkg/infrastructure/cache"
	ds "github.com/antsrp/banner_service/pkg/infrastructure/db"
	rs "github.com/antsrp/banner_service/pkg/infrastructure/rest"
	"github.com/antsrp/banner_service/pkg/jwt"
	"github.com/antsrp/banner_service/pkg/logger"
	"github.com/antsrp/banner_service/pkg/logger/slog"
)

func main() {
	var logger logger.Logger = slog.NewTextLogger(os.Stdout, slog.WithDebugLevel(), slog.WithFormat())
	key, err := os.ReadFile(".secret")
	if err != nil {
		logger.Fatal("can't parse secret key: %v", err.Error())
		return
	}
	if err := config.Load(); err != nil {
		logger.Fatal("can't load env files: %v", err.Error())
		return
	}

	dbSettings, err := config.Parse[ds.Settings]("DB")
	if err != nil {
		logger.Fatal("can't parse database settings from env file: %v", err.Error())
		return
	}
	dbConn, err := postgres.NewConnection(context.Background(), dbSettings, logger)
	if err != nil {
		logger.Fatal("can't create database connection: %v", err.Error())
		return
	}
	defer dbConn.Close()
	ustorage := postgres.NewUserStorage(dbConn)
	bstorage := postgres.NewBannerStorage(dbConn)

	serverSettings, err := config.Parse[rs.Settings]("SERVER")
	if err != nil {
		logger.Fatal("can't parse http server settings from env file: %v", err.Error())
	}
	js := jwt.NewJwtService(key, jwt.WithMethodHS256)

	cacheSettings, err := config.Parse[cs.Settings]("CACHE")
	if err != nil {
		logger.Fatal("can't parse cache settings from env file: %v", err.Error())
	}
	cacheStorage, err := redis.NewStorage[models.Banner](cacheSettings, logger)
	if err != nil {
		logger.Fatal("can't create redis connection: %v", err.Error())
	}

	bs := service.NewBannerService(bstorage, cacheStorage, logger)
	us := service.NewUserService(ustorage, js, logger)

	handler := rest.NewHandler(serverSettings, logger, bs, us)

	quit := make(chan struct{})
	transmitter := service.NewTransmitService(bs, cacheStorage, logger, quit)
	go transmitter.Start()

	if err := handler.Run(); err != nil {
		logger.Error("can't run http server: %v", err.Error())
	}

	transmitter.Stop()
}
