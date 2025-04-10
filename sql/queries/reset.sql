-- name: DeleteAllUsers :one
DELETE FROM users
RETURNING *;