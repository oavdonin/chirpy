-- name: DeleteAllChirpsFromUser :exec
DELETE FROM chirps WHERE user_id = $1;
DELETE FROM users WHERE id = $1;
