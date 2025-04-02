-- +goose Up
CREATE TABLE users (
    id uuid PRIMARY KEY,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
    email varchar(255) UNIQUE
);
-- +goose Down
DROP TABLE users;