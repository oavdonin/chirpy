-- +goose Up
CREATE TABLE chirps (
    id uuid PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    body TEXT,
    user_id uuid REFERENCES users(id) ON DELETE CASCADE
);
-- +goose Down
DROP TABLE chirps;