package models

import "time"

type BannerContent map[string]any

type BannerCommon struct {
	ID        int           `json:"id,omitempty"`
	FeatureID int           `json:"feature_id,omitempty"`
	TagIDS    []int         `json:"tag_ids,omitempty"`
	Content   BannerContent `json:"content,omitempty"`
	IsActive  *bool         `json:"is_active,omitempty"`
}

type Banner struct {
	BannerCommon
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
