package repos

import (
	"context"
	"time"

	"foodflow/internal/core"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClaimRepo interface {
	CreateClaim(ctx context.Context, orgID, offerID string, requestedServings int, status core.ClaimStatus) (*core.DonationClaim, error)
	GetClaimByID(ctx context.Context, id string) (*core.DonationClaim, error)
	GetClaimsByOrganization(ctx context.Context, orgID string, limit, offset int) ([]*core.DonationClaim, error)
	GetClaimsByStatus(ctx context.Context, status core.ClaimStatus, limit, offset int) ([]*core.DonationClaim, error)
	UpdateClaimStatus(ctx context.Context, id string, status core.ClaimStatus) (*core.DonationClaim, error)
	GetAllClaims(ctx context.Context, limit, offset int) ([]*core.DonationClaim, error)
}

type claimRepo struct {
	pool *pgxpool.Pool
}

func NewClaimRepo(pool *pgxpool.Pool) ClaimRepo {
	return &claimRepo{
		pool: pool,
	}
}

func (r *claimRepo) CreateClaim(ctx context.Context, orgID, offerID string, requestedServings int, status core.ClaimStatus) (*core.DonationClaim, error) {
	id := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO donation_claims (id, org_id, offer_id, requested_servings, priority_score, status, created_at)
		VALUES ($1, $2, $3, $4, 0.0, $5, $6)
		RETURNING id, offer_id, org_id, requested_servings, priority_score, status, created_at
	`

	var claim core.DonationClaim
	err := r.pool.QueryRow(ctx, query, id, orgID, offerID, requestedServings, string(status), now).Scan(
		&claim.ID, &claim.OfferID, &claim.OrgID, &claim.RequestedServings, &claim.PriorityScore, &claim.Status, &claim.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &claim, nil
}

func (r *claimRepo) GetClaimByID(ctx context.Context, id string) (*core.DonationClaim, error) {
	query := `
		SELECT id, offer_id, org_id, requested_servings, priority_score, status, created_at
		FROM donation_claims
		WHERE id = $1
	`

	var claim core.DonationClaim
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&claim.ID, &claim.OfferID, &claim.OrgID, &claim.RequestedServings, &claim.PriorityScore, &claim.Status, &claim.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, core.NewNotFoundError(core.ErrMsgClaimNotFound)
		}
		return nil, err
	}

	return &claim, nil
}

func (r *claimRepo) GetAllClaims(ctx context.Context, limit, offset int) ([]*core.DonationClaim, error) {
	query := `
		SELECT id, offer_id, org_id, requested_servings, priority_score, status, created_at
		FROM donation_claims
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var claims []*core.DonationClaim
	for rows.Next() {
		var claim core.DonationClaim
		err := rows.Scan(
			&claim.ID, &claim.OfferID, &claim.OrgID, &claim.RequestedServings, &claim.PriorityScore, &claim.Status, &claim.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		claims = append(claims, &claim)
	}

	return claims, nil
}

func (r *claimRepo) GetClaimsByOrganization(ctx context.Context, orgID string, limit, offset int) ([]*core.DonationClaim, error) {
	query := `
		SELECT id, offer_id, org_id, requested_servings, priority_score, status, created_at
		FROM donation_claims
		WHERE org_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var claims []*core.DonationClaim
	for rows.Next() {
		var claim core.DonationClaim
		err := rows.Scan(
			&claim.ID, &claim.OfferID, &claim.OrgID, &claim.RequestedServings, &claim.PriorityScore, &claim.Status, &claim.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		claims = append(claims, &claim)
	}

	return claims, nil
}

func (r *claimRepo) GetClaimsByStatus(ctx context.Context, status core.ClaimStatus, limit, offset int) ([]*core.DonationClaim, error) {
	query := `
		SELECT id, offer_id, org_id, requested_servings, priority_score, status, created_at
		FROM donation_claims
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, string(status), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var claims []*core.DonationClaim
	for rows.Next() {
		var claim core.DonationClaim
		err := rows.Scan(
			&claim.ID, &claim.OfferID, &claim.OrgID, &claim.RequestedServings, &claim.PriorityScore, &claim.Status, &claim.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		claims = append(claims, &claim)
	}

	return claims, nil
}

func (r *claimRepo) UpdateClaimStatus(ctx context.Context, id string, status core.ClaimStatus) (*core.DonationClaim, error) {
	query := `
		UPDATE donation_claims 
		SET status = $1
		WHERE id = $2
		RETURNING id, offer_id, org_id, requested_servings, priority_score, status, created_at
	`

	var claim core.DonationClaim
	err := r.pool.QueryRow(ctx, query, string(status), id).Scan(
		&claim.ID, &claim.OfferID, &claim.OrgID, &claim.RequestedServings, &claim.PriorityScore, &claim.Status, &claim.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, core.NewNotFoundError(core.ErrMsgClaimNotFound)
		}
		return nil, err
	}

	return &claim, nil
}