-- +goose Up
-- +goose StatementBegin
CREATE TABLE banners (
    id SERIAL PRIMARY KEY,
    is_active BOOLEAN DEFAULT TRUE,
    content JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), 
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    feature_id REFERENCES features(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE banners;
-- +goose StatementEnd
