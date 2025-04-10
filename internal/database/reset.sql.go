// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: reset.sql

package database

import (
	"context"
)

const deleteAllUsers = `-- name: DeleteAllUsers :one
DELETE FROM users
RETURNING id, created_at, updated_at, email
`

func (q *Queries) DeleteAllUsers(ctx context.Context) (User, error) {
	row := q.db.QueryRowContext(ctx, deleteAllUsers)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
	)
	return i, err
}
