package repos

import (
	"context"
	"time"

	"foodflow/internal/core"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrgRepo interface {
	CreateOrganization(ctx context.Context, userID string, orgType string, purposeFocus []string, lastMonthPeopleFed int) (*core.Organization, error)
	GetOrganizationByID(ctx context.Context, id string) (*core.Organization, error)
	GetOrganizationByUserID(ctx context.Context, userID string) (*core.Organization, error)
	UpdateOrganization(ctx context.Context, id string, orgType string, purposeFocus []string, lastMonthPeopleFed int) (*core.Organization, error)
	ListOrganizations(ctx context.Context, limit, offset int) ([]*core.Organization, error)
	GetAllOrganizations(ctx context.Context, limit, offset int) ([]*core.Organization, error)
}

type orgRepo struct {
	pool *pgxpool.Pool
}

func NewOrgRepo(pool *pgxpool.Pool) OrgRepo {
	return &orgRepo{
		pool: pool,
	}
}

func (r *orgRepo) CreateOrganization(ctx context.Context, userID string, orgType string, purposeFocus []string, lastMonthPeopleFed int) (*core.Organization, error) {
	id := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO organizations (id, user_id, org_type, purpose_focus, last_month_people_fed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, org_type, purpose_focus, last_month_people_fed, created_at, updated_at
	`

	var org core.Organization
	err := r.pool.QueryRow(ctx, query, id, userID, orgType, purposeFocus, lastMonthPeopleFed, now, now).Scan(
		&org.ID, &org.UserID, &org.OrgType, &org.PurposeFocus, &org.LastMonthPeopleFed, &org.CreatedAt, &org.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &org, nil
}

func (r *orgRepo) GetOrganizationByID(ctx context.Context, id string) (*core.Organization, error) {
	query := `
		SELECT id, user_id, org_type, purpose_focus, last_month_people_fed, created_at, updated_at
		FROM organizations
		WHERE id = $1
	`

	var org core.Organization
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&org.ID, &org.UserID, &org.OrgType, &org.PurposeFocus, &org.LastMonthPeopleFed, &org.CreatedAt, &org.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, core.NewNotFoundError(core.ErrMsgOrganizationNotFound)
		}
		return nil, err
	}

	return &org, nil
}

func (r *orgRepo) GetOrganizationByUserID(ctx context.Context, userID string) (*core.Organization, error) {
	query := `
		SELECT id, user_id, org_type, purpose_focus, last_month_people_fed, created_at, updated_at
		FROM organizations
		WHERE user_id = $1
	`

	var org core.Organization
	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&org.ID, &org.UserID, &org.OrgType, &org.PurposeFocus, &org.LastMonthPeopleFed, &org.CreatedAt, &org.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, core.NewNotFoundError(core.ErrMsgOrganizationNotFound)
		}
		return nil, err
	}

	return &org, nil
}

func (r *orgRepo) UpdateOrganization(ctx context.Context, id string, orgType string, purposeFocus []string, lastMonthPeopleFed int) (*core.Organization, error) {
	query := `
		UPDATE organizations 
		SET org_type = $1, purpose_focus = $2, last_month_people_fed = $3, updated_at = $4
		WHERE id = $5
		RETURNING id, user_id, org_type, purpose_focus, last_month_people_fed, created_at, updated_at
	`

	var org core.Organization
	err := r.pool.QueryRow(ctx, query, orgType, purposeFocus, lastMonthPeopleFed, time.Now(), id).Scan(
		&org.ID, &org.UserID, &org.OrgType, &org.PurposeFocus, &org.LastMonthPeopleFed, &org.CreatedAt, &org.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, core.NewNotFoundError(core.ErrMsgOrganizationNotFound)
		}
		return nil, err
	}

	return &org, nil
}

func (r *orgRepo) ListOrganizations(ctx context.Context, limit, offset int) ([]*core.Organization, error) {
	query := `
		SELECT id, user_id, org_type, purpose_focus, last_month_people_fed, created_at, updated_at
		FROM organizations
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orgs []*core.Organization
	for rows.Next() {
		var org core.Organization
		err := rows.Scan(
			&org.ID, &org.UserID, &org.OrgType, &org.PurposeFocus, &org.LastMonthPeopleFed, &org.CreatedAt, &org.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		orgs = append(orgs, &org)
	}

	return orgs, nil
}

func (r *orgRepo) GetAllOrganizations(ctx context.Context, limit, offset int) ([]*core.Organization, error) {
	return r.ListOrganizations(ctx, limit, offset)
}
