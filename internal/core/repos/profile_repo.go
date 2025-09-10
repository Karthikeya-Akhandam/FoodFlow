package repos

import (
	"context"
	"time"

	"foodflow/internal/core"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProfileRepo interface {
	CreateProfile(ctx context.Context, userID string, name string, phone, pincode, city, district, state, address *string, isRemoteOrg bool) (*core.Profile, error)
	GetProfileByUserID(ctx context.Context, userID string) (*core.Profile, error)
	UpdateProfile(ctx context.Context, userID string, name string, phone, pincode, city, district, state, address *string, isRemoteOrg bool) (*core.Profile, error)
}

type profileRepo struct {
	pool *pgxpool.Pool
}

func NewProfileRepo(pool *pgxpool.Pool) ProfileRepo {
	return &profileRepo{
		pool: pool,
	}
}

func (r *profileRepo) CreateProfile(ctx context.Context, userID string, name string, phone, pincode, city, district, state, address *string, isRemoteOrg bool) (*core.Profile, error) {
	now := time.Now()

	query := `
		INSERT INTO profiles (user_id, name, phone, pincode, city, district, state, address, is_remote_org, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING user_id, name, phone, pincode, city, district, state, address, is_remote_org, created_at, updated_at
	`

	var profile core.Profile
	err := r.pool.QueryRow(ctx, query, userID, name, phone, pincode, city, district, state, address, isRemoteOrg, now, now).Scan(
		&profile.UserID, &profile.Name, &profile.Phone, &profile.Pincode, &profile.City, &profile.District, &profile.State, &profile.Address, &profile.IsRemoteOrg, &profile.CreatedAt, &profile.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (r *profileRepo) GetProfileByUserID(ctx context.Context, userID string) (*core.Profile, error) {
	query := `
		SELECT user_id, name, phone, pincode, city, district, state, address, is_remote_org, created_at, updated_at
		FROM profiles
		WHERE user_id = $1
	`

	var profile core.Profile
	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&profile.UserID, &profile.Name, &profile.Phone, &profile.Pincode, &profile.City, &profile.District, &profile.State, &profile.Address, &profile.IsRemoteOrg, &profile.CreatedAt, &profile.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, core.NewNotFoundError(core.ErrMsgProfileNotFound)
		}
		return nil, err
	}

	return &profile, nil
}

func (r *profileRepo) UpdateProfile(ctx context.Context, userID string, name string, phone, pincode, city, district, state, address *string, isRemoteOrg bool) (*core.Profile, error) {
	query := `
		UPDATE profiles 
		SET name = $1, phone = $2, pincode = $3, city = $4, district = $5, state = $6, address = $7, is_remote_org = $8, updated_at = $9
		WHERE user_id = $10
		RETURNING user_id, name, phone, pincode, city, district, state, address, is_remote_org, created_at, updated_at
	`

	var profile core.Profile
	err := r.pool.QueryRow(ctx, query, name, phone, pincode, city, district, state, address, isRemoteOrg, time.Now(), userID).Scan(
		&profile.UserID, &profile.Name, &profile.Phone, &profile.Pincode, &profile.City, &profile.District, &profile.State, &profile.Address, &profile.IsRemoteOrg, &profile.CreatedAt, &profile.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, core.NewNotFoundError(core.ErrMsgProfileNotFound)
		}
		return nil, err
	}

	return &profile, nil
}
