CREATE TABLE IF NOT EXISTS urls (
    short_url UUID PRIMARY KEY,
    original_url VARCHAR(1000) NOT NULL
);