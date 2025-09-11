package repos

import (
	"context"
	"time"

	"foodflow/internal/core"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RedemptionRepo interface {
	CreateRedemption(ctx context.Context, offerID, orgID, collabID string, servingsAccepted, creditsSpent, tokensAwarded int) (*core.Redemption, error)
	GetRedemptionByID(ctx context.Context, id string) (*core.Redemption, error)
	GetRedemptionsByOrganization(ctx context.Context, orgID string, limit, offset int) ([]*core.Redemption, error)
	GetRedemptionsByCollaborator(ctx context.Context, collabID string, limit, offset int) ([]*core.Redemption, error)
	GetRedemptionsByOffer(ctx context.Context, offerID string, limit, offset int) ([]*core.Redemption, error)
	UpdateRedemption(ctx context.Context, id string, servingsAccepted, creditsSpent, tokensAwarded int) (*core.Redemption, error)
	GetAllRedemptions(ctx context.Context, limit, offset int) ([]*core.Redemption, error)
}

type redemptionRepo struct {
	pool *pgxpool.Pool
}

func NewRedemptionRepo(pool *pgxpool.Pool) RedemptionRepo {
	return &redemptionRepo{
		pool: pool,
	}
}

func (r *redemptionRepo) CreateRedemption(ctx context.Context, offerID, orgID, collabID string, servingsAccepted, creditsSpent, tokensAwarded int) (*core.Redemption, error) {
	id := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO redemptions (id, offer_id, org_id, collab_id, servings_accepted, credits_spent, tokens_awarded, confirmed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, offer_id, org_id, collab_id, servings_accepted, credits_spent, tokens_awarded, confirmed_at
	`

	var redemption core.Redemption
	err := r.pool.QueryRow(ctx, query, id, offerID, orgID, collabID, servingsAccepted, creditsSpent, tokensAwarded, now).Scan(
		&redemption.ID, &redemption.OfferID, &redemption.OrgID, &redemption.CollabID, &redemption.ServingsAccepted, &redemption.CreditsSpent, &redemption.TokensAwarded, &redemption.ConfirmedAt,
	)
	if err != nil {
		return nil, err
	}

	return &redemption, nil
}

func (r *redemptionRepo) GetRedemptionByID(ctx context.Context, id string) (*core.Redemption, error) {
	query := `
		SELECT id, offer_id, org_id, collab_id, servings_accepted, credits_spent, tokens_awarded, confirmed_at
		FROM redemptions
		WHERE id = $1
	`

	var redemption core.Redemption
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&redemption.ID, &redemption.OfferID, &redemption.OrgID, &redemption.CollabID, &redemption.ServingsAccepted, &redemption.CreditsSpent, &redemption.TokensAwarded, &redemption.ConfirmedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, core.NewNotFoundError(core.ErrMsgRedemptionNotFound)
		}
		return nil, err
	}

	return &redemption, nil
}

func (r *redemptionRepo) GetRedemptionsByOrganization(ctx context.Context, orgID string, limit, offset int) ([]*core.Redemption, error) {
	query := `
		SELECT id, offer_id, org_id, collab_id, servings_accepted, credits_spent, tokens_awarded, confirmed_at
		FROM redemptions
		WHERE org_id = $1
		ORDER BY confirmed_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var redemptions []*core.Redemption
	for rows.Next() {
		var redemption core.Redemption
		err := rows.Scan(
			&redemption.ID, &redemption.OfferID, &redemption.OrgID, &redemption.CollabID, &redemption.ServingsAccepted, &redemption.CreditsSpent, &redemption.TokensAwarded, &redemption.ConfirmedAt,
		)
		if err != nil {
			return nil, err
		}
		redemptions = append(redemptions, &redemption)
	}

	return redemptions, nil
}

func (r *redemptionRepo) GetRedemptionsByOffer(ctx context.Context, offerID string, limit, offset int) ([]*core.Redemption, error) {
	query := `
		SELECT id, offer_id, org_id, collab_id, servings_accepted, credits_spent, tokens_awarded, confirmed_at
		FROM redemptions
		WHERE offer_id = $1
		ORDER BY confirmed_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, offerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var redemptions []*core.Redemption
	for rows.Next() {
		var redemption core.Redemption
		err := rows.Scan(
			&redemption.ID, &redemption.OfferID, &redemption.OrgID, &redemption.CollabID, &redemption.ServingsAccepted, &redemption.CreditsSpent, &redemption.TokensAwarded, &redemption.ConfirmedAt,
		)
		if err != nil {
			return nil, err
		}
		redemptions = append(redemptions, &redemption)
	}

	return redemptions, nil
}

func (r *redemptionRepo) GetRedemptionsByCollaborator(ctx context.Context, collabID string, limit, offset int) ([]*core.Redemption, error) {
	query := `
		SELECT id, offer_id, org_id, collab_id, servings_accepted, credits_spent, tokens_awarded, confirmed_at
		FROM redemptions
		WHERE collab_id = $1
		ORDER BY confirmed_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, collabID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var redemptions []*core.Redemption
	for rows.Next() {
		var redemption core.Redemption
		err := rows.Scan(
			&redemption.ID, &redemption.OfferID, &redemption.OrgID, &redemption.CollabID, &redemption.ServingsAccepted, &redemption.CreditsSpent, &redemption.TokensAwarded, &redemption.ConfirmedAt,
		)
		if err != nil {
			return nil, err
		}
		redemptions = append(redemptions, &redemption)
	}

	return redemptions, nil
}

func (r *redemptionRepo) UpdateRedemption(ctx context.Context, id string, servingsAccepted, creditsSpent, tokensAwarded int) (*core.Redemption, error) {
	query := `
		UPDATE redemptions 
		SET servings_accepted = $1, credits_spent = $2, tokens_awarded = $3
		WHERE id = $4
		RETURNING id, offer_id, org_id, collab_id, servings_accepted, credits_spent, tokens_awarded, confirmed_at
	`

	var redemption core.Redemption
	err := r.pool.QueryRow(ctx, query, servingsAccepted, creditsSpent, tokensAwarded, id).Scan(
		&redemption.ID, &redemption.OfferID, &redemption.OrgID, &redemption.CollabID, &redemption.ServingsAccepted, &redemption.CreditsSpent, &redemption.TokensAwarded, &redemption.ConfirmedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, core.NewNotFoundError(core.ErrMsgRedemptionNotFound)
		}
		return nil, err
	}

	return &redemption, nil
}

func (r *redemptionRepo) GetAllRedemptions(ctx context.Context, limit, offset int) ([]*core.Redemption, error) {
	query := `
		SELECT id, offer_id, org_id, collab_id, servings_accepted, credits_spent, tokens_awarded, confirmed_at
		FROM redemptions
		ORDER BY confirmed_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var redemptions []*core.Redemption
	for rows.Next() {
		var redemption core.Redemption
		err := rows.Scan(
			&redemption.ID, &redemption.OfferID, &redemption.OrgID, &redemption.CollabID, &redemption.ServingsAccepted, &redemption.CreditsSpent, &redemption.TokensAwarded, &redemption.ConfirmedAt,
		)
		if err != nil {
			return nil, err
		}
		redemptions = append(redemptions, &redemption)
	}

	return redemptions, nil
}
