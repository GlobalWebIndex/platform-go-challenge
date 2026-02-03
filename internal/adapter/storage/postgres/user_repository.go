package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"go-favorites-app/internal/core/domain/auth"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Save(ctx context.Context, user auth.User) error {
	query := `INSERT INTO users (id, email, password_hash) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, user.ID, user.Email, user.PasswordHash)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}
	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (auth.User, error) {
	query := `SELECT id, email, password_hash FROM users WHERE email = $1`
	row := r.db.QueryRow(ctx, query, email)

	var user auth.User
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return auth.User{}, fmt.Errorf("failed to find user: %w", err)
	}
	return user, nil
}
