CREATE TABLE banners (
    id SERIAL PRIMARY KEY,
    is_active BOOLEAN DEFAULT TRUE,
    content JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), 
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    feature_id INTEGER REFERENCES features(id)
);