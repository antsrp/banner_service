CREATE TABLE banners_tags (
    id SERIAL PRIMARY KEY,
    banner_id INTEGER REFERENCES banners(id) ON DELETE CASCADE ON UPDATE CASCADE,
    tag_id INTEGER REFERENCES tags(id)
);