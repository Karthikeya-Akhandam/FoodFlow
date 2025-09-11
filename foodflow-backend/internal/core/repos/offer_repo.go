package repos

import (
	"context"
	"time"

	"foodflow/internal/core"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OfferRepo interface {
	CreateOffer(ctx context.Context, collabID string, title string, description *string, readyFrom, expiresAt string, estimatedServings int, purpose, pincode, city, state *string, status core.DonationStatus) (*core.DonationOffer, error)
	GetOfferByID(ctx context.Context, id string) (*core.DonationOffer, error)
	GetOffersByCollaborator(ctx context.Context, collabID string, limit, offset int) ([]*core.DonationOffer, error)
	GetOffersByStatus(ctx context.Context, status core.DonationStatus, limit, offset int) ([]*core.DonationOffer, error)
	GetNearbyOffers(ctx context.Context, pincode, city, state string, purpose *string, limit, offset int) ([]*core.DonationOffer, error)
	UpdateOfferStatus(ctx context.Context, id string, status core.DonationStatus) (*core.DonationOffer, error)
	GetExpiredOffers(ctx context.Context) ([]*core.DonationOffer, error)
}

type offerRepo struct {
	pool *pgxpool.Pool
}

func NewOfferRepo(pool *pgxpool.Pool) OfferRepo {
	return &offerRepo{
		pool: pool,
	}
}

func (r *offerRepo) CreateOffer(ctx context.Context, collabID string, title string, description *string, readyFrom, expiresAt string, estimatedServings int, purpose, pincode, city, state *string, status core.DonationStatus) (*core.DonationOffer, error) {
	id := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO donation_offers (id, collab_id, title, description, ready_from, expires_at, estimated_servings, purpose, pincode, city, state, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, collab_id, title, description, ready_from, expires_at, estimated_servings, purpose, pincode, city, state, status, created_at, updated_at
	`

	var offer core.DonationOffer
	err := r.pool.QueryRow(ctx, query, id, collabID, title, description, readyFrom, expiresAt, estimatedServings, purpose, pincode, city, state, string(status), now, now).Scan(
		&offer.ID, &offer.CollabID, &offer.Title, &offer.Description, &offer.ReadyFrom, &offer.ExpiresAt, &offer.EstimatedServings, &offer.Purpose, &offer.Pincode, &offer.City, &offer.State, &offer.Status, &offer.CreatedAt, &offer.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &offer, nil
}

func (r *offerRepo) GetOfferByID(ctx context.Context, id string) (*core.DonationOffer, error) {
	query := `
		SELECT id, collab_id, title, description, ready_from, expires_at, estimated_servings, purpose, pincode, city, state, status, created_at, updated_at
		FROM donation_offers
		WHERE id = $1
	`

	var offer core.DonationOffer
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&offer.ID, &offer.CollabID, &offer.Title, &offer.Description, &offer.ReadyFrom, &offer.ExpiresAt, &offer.EstimatedServings, &offer.Purpose, &offer.Pincode, &offer.City, &offer.State, &offer.Status, &offer.CreatedAt, &offer.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, core.NewNotFoundError(core.ErrMsgOfferNotFound)
		}
		return nil, err
	}

	return &offer, nil
}

func (r *offerRepo) GetOffersByCollaborator(ctx context.Context, collabID string, limit, offset int) ([]*core.DonationOffer, error) {
	query := `
		SELECT id, collab_id, title, description, ready_from, expires_at, estimated_servings, purpose, pincode, city, state, status, created_at, updated_at
		FROM donation_offers
		WHERE collab_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, collabID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var offers []*core.DonationOffer
	for rows.Next() {
		var offer core.DonationOffer
		err := rows.Scan(
			&offer.ID, &offer.CollabID, &offer.Title, &offer.Description, &offer.ReadyFrom, &offer.ExpiresAt, &offer.EstimatedServings, &offer.Purpose, &offer.Pincode, &offer.City, &offer.State, &offer.Status, &offer.CreatedAt, &offer.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		offers = append(offers, &offer)
	}

	return offers, nil
}

func (r *offerRepo) GetOffersByStatus(ctx context.Context, status core.DonationStatus, limit, offset int) ([]*core.DonationOffer, error) {
	query := `
		SELECT id, collab_id, title, description, ready_from, expires_at, estimated_servings, purpose, pincode, city, state, status, created_at, updated_at
		FROM donation_offers
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, string(status), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var offers []*core.DonationOffer
	for rows.Next() {
		var offer core.DonationOffer
		err := rows.Scan(
			&offer.ID, &offer.CollabID, &offer.Title, &offer.Description, &offer.ReadyFrom, &offer.ExpiresAt, &offer.EstimatedServings, &offer.Purpose, &offer.Pincode, &offer.City, &offer.State, &offer.Status, &offer.CreatedAt, &offer.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		offers = append(offers, &offer)
	}

	return offers, nil
}

func (r *offerRepo) GetNearbyOffers(ctx context.Context, pincode, city, state string, purpose *string, limit, offset int) ([]*core.DonationOffer, error) {
	query := `
		SELECT id, collab_id, title, description, ready_from, expires_at, estimated_servings, purpose, pincode, city, state, status, created_at, updated_at
		FROM donation_offers
		WHERE status = 'available' 
		AND (pincode = $1 OR city = $2 OR state = $3)
		AND ($4 IS NULL OR purpose = $4)
		ORDER BY created_at DESC
		LIMIT $5 OFFSET $6
	`

	rows, err := r.pool.Query(ctx, query, pincode, city, state, purpose, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var offers []*core.DonationOffer
	for rows.Next() {
		var offer core.DonationOffer
		err := rows.Scan(
			&offer.ID, &offer.CollabID, &offer.Title, &offer.Description, &offer.ReadyFrom, &offer.ExpiresAt, &offer.EstimatedServings, &offer.Purpose, &offer.Pincode, &offer.City, &offer.State, &offer.Status, &offer.CreatedAt, &offer.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		offers = append(offers, &offer)
	}

	return offers, nil
}

func (r *offerRepo) UpdateOfferStatus(ctx context.Context, id string, status core.DonationStatus) (*core.DonationOffer, error) {
	query := `
		UPDATE donation_offers 
		SET status = $1, updated_at = $2
		WHERE id = $3
		RETURNING id, collab_id, title, description, ready_from, expires_at, estimated_servings, purpose, pincode, city, state, status, created_at, updated_at
	`

	var offer core.DonationOffer
	err := r.pool.QueryRow(ctx, query, string(status), time.Now(), id).Scan(
		&offer.ID, &offer.CollabID, &offer.Title, &offer.Description, &offer.ReadyFrom, &offer.ExpiresAt, &offer.EstimatedServings, &offer.Purpose, &offer.Pincode, &offer.City, &offer.State, &offer.Status, &offer.CreatedAt, &offer.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, core.NewNotFoundError(core.ErrMsgOfferNotFound)
		}
		return nil, err
	}

	return &offer, nil
}

func (r *offerRepo) GetExpiredOffers(ctx context.Context) ([]*core.DonationOffer, error) {
	query := `
		SELECT id, collab_id, title, description, ready_from, expires_at, estimated_servings, purpose, pincode, city, state, status, created_at, updated_at
		FROM donation_offers
		WHERE expires_at < NOW() AND status = 'available'
		ORDER BY expires_at ASC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var offers []*core.DonationOffer
	for rows.Next() {
		var offer core.DonationOffer
		err := rows.Scan(
			&offer.ID, &offer.CollabID, &offer.Title, &offer.Description, &offer.ReadyFrom, &offer.ExpiresAt, &offer.EstimatedServings, &offer.Purpose, &offer.Pincode, &offer.City, &offer.State, &offer.Status, &offer.CreatedAt, &offer.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		offers = append(offers, &offer)
	}

	return offers, nil
}
