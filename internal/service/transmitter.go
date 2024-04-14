package service

import (
	"fmt"
	"time"

	"github.com/antsrp/banner_service/internal/cache"
	"github.com/antsrp/banner_service/internal/domain/models"
	"github.com/antsrp/banner_service/internal/domain/models/requests"
	"github.com/antsrp/banner_service/pkg/logger"
)

type Transmitter interface {
	Start()
	Stop()
}

type TransmitService struct {
	bannerService BannerServicer
	cacheStorage  cache.Storager[models.Banner]
	logger        logger.Logger
	end           chan struct{}
}

func NewTransmitService(bs BannerServicer, cs cache.Storager[models.Banner], logger logger.Logger, end chan struct{}) TransmitService {
	return TransmitService{
		bannerService: bs,
		cacheStorage:  cs,
		logger:        logger,
		end:           end,
	}
}

func (s TransmitService) Start() {
	s.getAll()
	for {
		select {
		case <-s.end:
			return
		case <-time.After(270 * time.Second):
			s.getAll()
		}
	}
}

func (s TransmitService) Stop() {
	s.end <- struct{}{}
}

func (s TransmitService) getAll() {
	banners, err := s.bannerService.Get(requests.GetBannersRequest{})
	if err != nil {
		s.logger.Info("can't get banners to put them into cache: %v", err.Cause().Error())
	}
	if err := s.writeToCache(banners); err != nil {
		s.logger.Info("can't write banners into cache: %v", err.Error())
	}
}

func (s TransmitService) writeToCache(banners []models.Banner) error {
	for _, banner := range banners {
		for _, tag := range banner.TagIDS {
			key := fmt.Sprintf("(%d, %d)", banner.FeatureID, tag)
			if err := s.cacheStorage.Set(key, banner); err != nil {
				return fmt.Errorf("can't put banner into cache: %v", err.Error())
			}
		}
	}
	return nil
}

var _ Transmitter = TransmitService{}
