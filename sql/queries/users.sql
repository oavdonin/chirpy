-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    now(),
    now(),
    $1,
    $2
)
RETURNING *;
-- name: GetUserHash :one
SELECT hashed_password FROM users where email = $1;
-- name: GetUser :one
SELECT id, created_at, updated_at, email FROM users where email = $1;