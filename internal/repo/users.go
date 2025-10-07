package repo

import (
	"context"
	"database/sql"

	"github.com/SpBalis/platform-go-challenge/internal/service"
)

type UsersRepo struct{ DB *sql.DB }

func NewUsersRepo(db *sql.DB) *UsersRepo { return &UsersRepo{DB: db} }

func (r *UsersRepo) Create(ctx context.Context, email string) (int64, error) {
	const q = `INSERT INTO users (email) VALUES ($1) RETURNING id`
	var id int64
	if err := r.DB.QueryRowContext(ctx, q, email).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UsersRepo) List(ctx context.Context, limit, offset int) ([]service.User, error) {
	const q = `SELECT id, email FROM users ORDER BY id LIMIT $1 OFFSET $2`
	rows, err := r.DB.QueryContext(ctx, q, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []service.User
	for rows.Next() {
		var u service.User
		if err := rows.Scan(&u.ID, &u.Email); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

func (r *UsersRepo) Get(ctx context.Context, id int64) (service.User, error) {
	const q = `SELECT id, email FROM users WHERE id=$1`
	var u service.User
	err := r.DB.QueryRowContext(ctx, q, id).Scan(&u.ID, &u.Email)
	return u, err
}
