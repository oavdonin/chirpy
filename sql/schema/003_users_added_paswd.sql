-- +goose Up
ALTER TABLE users ADD COLUMN status TEXT DEFAULT 'unset';

UPDATE users SET status = 'unset' WHERE status IS NULL;

ALTER TABLE users ALTER COLUMN status SET NOT NULL;

-- +goose Down
ALTER TABLE users DROP COLUMN status;