package repository

import (
	"context"

	"github.com/antsrp/banner_service/internal/domain/models"
)

type GetBanner struct {
	FeatureID int
	TagID     int
}

type GetBannerLimited struct {
	GetBanner
	Limit  int
	Offset int
}

type BannerStorage interface {
	Create(context.Context, models.Banner) (models.Banner, DatabaseError)
	Update(context.Context, models.Banner) DatabaseError
	Get(ctx context.Context, opts GetBannerLimited) ([]models.Banner, DatabaseError)
	GetOne(ctx context.Context, opts GetBanner, userName string) (models.Banner, DatabaseError)
	Delete(ctx context.Context, id int) DatabaseError
}
