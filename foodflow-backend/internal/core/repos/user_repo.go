package repos

import (
	"context"
	"time"

	"foodflow/internal/core"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo interface {
	CreateUser(ctx context.Context, email, passwordHash string, role core.UserRole, status core.UserStatus) (*core.User, error)
	GetUserByID(ctx context.Context, id string) (*core.User, error)
	GetUserByEmail(ctx context.Context, email string) (*core.User, error)
	UpdateUserStatus(ctx context.Context, id string, status core.UserStatus) (*core.User, error)
	GetAllUsers(ctx context.Context, limit, offset int) ([]*core.User, error)
}

type userRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) UserRepo {
	return &userRepo{
		pool: pool,
	}
}

func (r *userRepo) CreateUser(ctx context.Context, email, passwordHash string, role core.UserRole, status core.UserStatus) (*core.User, error) {
	id := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO users (id, email, password_hash, role, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, email, password_hash, role, status, created_at, updated_at
	`

	var user core.User
	err := r.pool.QueryRow(ctx, query, id, email, passwordHash, string(role), string(status), now, now).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.Status, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) GetUserByID(ctx context.Context, id string) (*core.User, error) {
	query := `
		SELECT id, email, password_hash, role, status, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user core.User
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.Status, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, core.NewNotFoundError(core.ErrMsgUserNotFound)
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*core.User, error) {
	query := `
		SELECT id, email, password_hash, role, status, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user core.User
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.Status, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, core.NewNotFoundError(core.ErrMsgUserNotFound)
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) UpdateUserStatus(ctx context.Context, id string, status core.UserStatus) (*core.User, error) {
	query := `
		UPDATE users 
		SET status = $1, updated_at = $2
		WHERE id = $3
		RETURNING id, email, password_hash, role, status, created_at, updated_at
	`

	var user core.User
	err := r.pool.QueryRow(ctx, query, string(status), time.Now(), id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.Status, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, core.NewNotFoundError(core.ErrMsgUserNotFound)
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) GetAllUsers(ctx context.Context, limit, offset int) ([]*core.User, error) {
	query := `
		SELECT id, email, password_hash, role, status, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*core.User
	for rows.Next() {
		var user core.User
		err := rows.Scan(
			&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.Status, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}
