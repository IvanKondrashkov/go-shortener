CREATE TABLE IF NOT EXISTS urls (
    short_url UUID PRIMARY KEY,
    user_id UUID NULL,
    original_url VARCHAR(1000) NOT NULL
);