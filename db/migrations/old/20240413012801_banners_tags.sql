-- +goose Up
-- +goose StatementBegin
CREATE TABLE banners_tags (
    id SERIAL PRIMARY KEY,
    banner_id INTEGER REFERENCES banners(id),
    tag_id INTEGER REFERENCES tags(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE banners_tags;
-- +goose StatementEnd
