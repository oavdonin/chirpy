-- +goose Up
CREATE TABLE refresh_tokens (
    token varchar(64) PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id uuid REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP NULL,
    revoked_at TIMESTAMP NULL
);
-- +goose Down
DROP TABLE refresh_tokens;