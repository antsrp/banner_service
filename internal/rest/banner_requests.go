package rest

import "time"

type BannerContent map[any]any

type UserBannerRequest struct {
	TagID             int  `json:"tag_id"`
	FeatureID         int  `json:"feature_id"`
	IsUseLastRevision bool `json:"use_last_revision"`
	// user token also here
}

type UserBannerResponse struct {
	Content      BannerContent `json:"content,omitempty"` // undefined structure
	ErrorMessage string        `json:"error,omitempty"`
}

type GetBannersRequest struct {
	UserBannerRequest
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type BannerCommon struct {
	FeatureID int           `json:"feature_id,omitempty"`
	TagIDS    []int         `json:"tag_ids,omitempty"`
	Content   BannerContent `json:"content,omitempty"`
	IsActive  *bool         `json:"is_active,omitempty"`
}

type GetBannersResponse struct {
	BannerCommon
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
	ErrorMessage string    `json:"error,omitempty"`
}

type CreateBannerRequest struct {
	BannerCommon
}

type CreateBannerResponse struct {
	BannerID     int    `json:"banner_id,omitempty"`
	ErrorMessage string `json:"error,omitempty"`
}

type UpdateBannerRequest struct {
	ID int `json:"id"`
	BannerCommon
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
