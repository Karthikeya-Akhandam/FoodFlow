package repos

import (
	"context"
	"time"

	"foodflow/internal/core"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreditRepo interface {
	CreateCredit(ctx context.Context, orgID string, month string, creditsIssued int, isRemoteBonus bool) (*core.Credit, error)
	GetCreditsByOrganization(ctx context.Context, orgID string) ([]*core.Credit, error)
	GetCreditHistory(ctx context.Context, orgID string, limit, offset int) ([]*core.Credit, error)
	GetTotalCredits(ctx context.Context, orgID string) (int, error)
	GetCreditByOrgAndMonth(ctx context.Context, orgID string, month string) (*core.Credit, error)
	UpdateCreditRemaining(ctx context.Context, creditID string, creditsRemaining int) error
	SpendCredits(ctx context.Context, orgID string, amount int, reason string) (*core.Credit, error)
	IssueCredits(ctx context.Context, orgID string, month string, creditsToIssue int, isRemote bool) error
}

type creditRepo struct {
	pool *pgxpool.Pool
}

func NewCreditRepo(pool *pgxpool.Pool) CreditRepo {
	return &creditRepo{
		pool: pool,
	}
}

func (r *creditRepo) CreateCredit(ctx context.Context, orgID string, month string, creditsIssued int, isRemoteBonus bool) (*core.Credit, error) {
	id := uuid.New().String()
	now := time.Now()
	expiresAt := now.AddDate(0, 1, 0) // Expires after 1 month

	query := `
		INSERT INTO credits (id, org_id, month, credits_issued, credits_remaining, is_remote_bonus, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, org_id, month, credits_issued, credits_remaining, is_remote_bonus, expires_at, created_at
	`

	var credit core.Credit
	err := r.pool.QueryRow(ctx, query, id, orgID, month, creditsIssued, creditsIssued, isRemoteBonus, expiresAt, now).Scan(
		&credit.ID, &credit.OrgID, &credit.Month, &credit.CreditsIssued, &credit.CreditsRemaining, &credit.IsRemoteBonus, &credit.ExpiresAt, &credit.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &credit, nil
}

func (r *creditRepo) GetCreditsByOrganization(ctx context.Context, orgID string) ([]*core.Credit, error) {
	query := `
		SELECT id, org_id, month, credits_issued, credits_remaining, is_remote_bonus, expires_at, created_at
		FROM credits
		WHERE org_id = $1 AND credits_remaining > 0 AND expires_at > NOW()
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var credits []*core.Credit
	for rows.Next() {
		var credit core.Credit
		err := rows.Scan(
			&credit.ID, &credit.OrgID, &credit.Month, &credit.CreditsIssued, &credit.CreditsRemaining, &credit.IsRemoteBonus, &credit.ExpiresAt, &credit.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		credits = append(credits, &credit)
	}

	return credits, nil
}

func (r *creditRepo) GetCreditHistory(ctx context.Context, orgID string, limit, offset int) ([]*core.Credit, error) {
	query := `
		SELECT id, org_id, month, credits_issued, credits_remaining, is_remote_bonus, expires_at, created_at
		FROM credits
		WHERE org_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var credits []*core.Credit
	for rows.Next() {
		var credit core.Credit
		err := rows.Scan(
			&credit.ID, &credit.OrgID, &credit.Month, &credit.CreditsIssued, &credit.CreditsRemaining, &credit.IsRemoteBonus, &credit.ExpiresAt, &credit.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		credits = append(credits, &credit)
	}

	return credits, nil
}

func (r *creditRepo) GetTotalCredits(ctx context.Context, orgID string) (int, error) {
	query := `
		SELECT COALESCE(SUM(credits_remaining), 0) as total
		FROM credits
		WHERE org_id = $1 AND expires_at > NOW()
	`

	var total int
	err := r.pool.QueryRow(ctx, query, orgID).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *creditRepo) GetCreditByOrgAndMonth(ctx context.Context, orgID string, month string) (*core.Credit, error) {
	query := `
		SELECT id, org_id, month, credits_issued, credits_remaining, is_remote_bonus, expires_at, created_at
		FROM credits
		WHERE org_id = $1 AND month = $2
		ORDER BY created_at DESC
		LIMIT 1
	`

	var credit core.Credit
	err := r.pool.QueryRow(ctx, query, orgID, month).Scan(
		&credit.ID, &credit.OrgID, &credit.Month, &credit.CreditsIssued, &credit.CreditsRemaining, &credit.IsRemoteBonus, &credit.ExpiresAt, &credit.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, core.NewNotFoundError(core.ErrMsgCreditNotFound)
		}
		return nil, err
	}

	return &credit, nil
}

func (r *creditRepo) UpdateCreditRemaining(ctx context.Context, creditID string, creditsRemaining int) error {
	query := `
		UPDATE credits 
		SET credits_remaining = $1
		WHERE id = $2
	`

	_, err := r.pool.Exec(ctx, query, creditsRemaining, creditID)
	if err != nil {
		return err
	}

	return nil
}

func (r *creditRepo) SpendCredits(ctx context.Context, orgID string, amount int, reason string) (*core.Credit, error) {
	// For now, we'll just return a placeholder credit
	// In a real implementation, this would create a credit transaction record
	return &core.Credit{
		ID:               uuid.New(),
		OrgID:            uuid.MustParse(orgID),
		Month:            time.Now().Format("200601"),
		CreditsIssued:    0,
		CreditsRemaining: -amount, // Negative to represent spending
		IsRemoteBonus:    false,
		ExpiresAt:        time.Now().AddDate(0, 1, 0),
		CreatedAt:        time.Now(),
	}, nil
}

func (r *creditRepo) IssueCredits(ctx context.Context, orgID string, month string, creditsToIssue int, isRemote bool) error {
	// Check if credits already exist for this org and month
	existing, err := r.GetCreditByOrgAndMonth(ctx, orgID, month)
	if err != nil && err != pgx.ErrNoRows {
		return err
	}
	
	if existing != nil {
		// Update existing credits
		newRemaining := existing.CreditsRemaining + creditsToIssue
		return r.UpdateCreditRemaining(ctx, existing.ID.String(), newRemaining)
	}
	
	// Create new credit record
	_, err = r.CreateCredit(ctx, orgID, month, creditsToIssue, isRemote)
	return err
}
