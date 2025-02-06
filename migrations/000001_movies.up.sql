CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE IF NOT EXISTS movies (
    id UUID PRIMARY KEY,
    name_uz VARCHAR(256) NOT NULL,
    name_en VARCHAR(256) NOT NULL,
    name_ru VARCHAR(256) NOT NULL,
    embedding VECTOR(768),
    created_at timestamp NOT NULL DEFAULT 'now()',
    updated_at timestamp NOT NULL DEFAULT 'now()'
);