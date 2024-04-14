package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/antsrp/banner_service/internal/cache"
	"github.com/antsrp/banner_service/internal/domain/models"
	"github.com/antsrp/banner_service/internal/domain/models/requests"
	"github.com/antsrp/banner_service/internal/repository"
	"github.com/antsrp/banner_service/internal/repository/postgres"
	"github.com/antsrp/banner_service/pkg/logger"
)

type BannerServicer interface {
	GetOne(requests.UserBannerRequest, string) (models.Banner, Error)
	Get(requests.GetBannersRequest) ([]models.Banner, Error)
	Create(requests.CreateBannerRequest) (models.Banner, Error)
	Update(requests.UpdateBannerRequest) Error
	Delete(requests.DeleteBannerRequest) Error
}

type BannerService struct {
	storage      postgres.BannerStorage
	cacheStorage cache.Storager[models.Banner]
	logger       logger.Logger
}

func NewBannerService(storage postgres.BannerStorage, cs cache.Storager[models.Banner], logger logger.Logger) BannerService {
	return BannerService{
		storage:      storage,
		cacheStorage: cs,
		logger:       logger,
	}
}

func (s BannerService) GetOne(req requests.UserBannerRequest, userName string) (models.Banner, Error) {
	var banner models.Banner
	if req.IsUseLastRevision { // find in cache
		var err error
		banner, err = s.cacheStorage.Get(fmt.Sprintf("(%d, %d)", req.FeatureID, req.TagID))
		if err != nil {
			return models.Banner{}, NewServiceError(false, ErrBannerNotFound)
		}
	} else {
		var err Error
		banner, err = s.storage.GetOne(context.Background(), repository.GetBanner{FeatureID: req.FeatureID, TagID: req.TagID}, userName)
		if err != nil {
			s.logger.Error("can't get banner: %v", err.Cause().Error())
			if errors.Is(err.Cause(), repository.ErrEntityNotFound) {
				return models.Banner{}, NewServiceError(false, ErrBannerNotFound)
			}
			return models.Banner{}, err
		}
	}
	return banner, nil
}
func (s BannerService) Get(req requests.GetBannersRequest) ([]models.Banner, Error) {
	banners, err := s.storage.Get(context.Background(), repository.GetBannerLimited{
		GetBanner: repository.GetBanner{
			FeatureID: req.FeatureID,
			TagID:     req.TagID,
		},
		Limit:  req.Limit,
		Offset: req.Offset,
	})
	if err != nil {
		if err.IsInternal() {
			return nil, defaultInternalError
		}
		return nil, NewServiceError(true, err.Cause())
	}
	return banners, nil
}
func (s BannerService) Create(req requests.CreateBannerRequest) (models.Banner, Error) {
	banner, err := s.storage.Create(context.Background(), models.Banner{
		BannerCommon: req.BannerCommon,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})
	if err != nil {
		if err.IsInternal() {
			return models.Banner{}, defaultInternalError
		}
		return models.Banner{}, NewServiceError(true, err.Cause())
	}
	return banner, nil
}
func (s BannerService) Update(req requests.UpdateBannerRequest) Error {
	if err := s.storage.Update(context.Background(), models.Banner{
		BannerCommon: req.BannerCommon,
		UpdatedAt:    time.Now(),
	}); err != nil {
		if errors.Is(err.Cause(), repository.ErrNoRowsAffected) {
			return NewServiceError(false, ErrBannerNotFound)
		}
		if err.IsInternal() {
			return defaultInternalError
		}
		return NewServiceError(true, err.Cause())
	}
	return nil
}
func (s BannerService) Delete(req requests.DeleteBannerRequest) Error {
	if err := s.storage.Delete(context.Background(), req.ID); err != nil {
		if errors.Is(err.Cause(), repository.ErrNoRowsAffected) {
			return NewServiceError(false, ErrBannerNotFound)
		}
		if err.IsInternal() {
			return defaultInternalError
		}
		return NewServiceError(true, err.Cause())
	}
	return nil
}

var _ BannerServicer = BannerService{}
