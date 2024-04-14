package requests

import (
	"github.com/antsrp/banner_service/internal/domain/models"
)

type UserBannerRequest struct {
	TagID             int  `json:"tag_id"`
	FeatureID         int  `json:"feature_id"`
	IsUseLastRevision bool `json:"use_last_revision"`
	// user token also here
}

type UserBannerResponse struct {
	Content      models.BannerContent `json:"content,omitempty"` // undefined structure
	ErrorMessage string               `json:"error,omitempty"`
}

type GetBannersRequest struct {
	UserBannerRequest
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type GetBannersResponse struct {
	models.Banner
	ErrorMessage string `json:"error,omitempty"`
}

type CreateBannerRequest struct {
	models.BannerCommon
}

type CreateBannerResponse struct {
	BannerID     int    `json:"banner_id,omitempty"`
	ErrorMessage string `json:"error,omitempty"`
}

type UpdateBannerRequest struct {
	models.BannerCommon
}

type UpdateBannerResponse struct {
	ErrorMessage string `json:"error,omitempty"`
}

type DeleteBannerRequest struct {
	ID int `json:"id"`
}

type DeleteBannerResponse struct {
	ErrorMessage string `json:"error,omitempty"`
}
