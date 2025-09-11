package repos

import (
	"context"
	"time"

	"foodflow/internal/core"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RemoteRepo interface {
	CreateRemoteOrgAssignment(ctx context.Context, remoteOrgID, proxyOrgID string, notes *string) (*core.RemoteOrgAssignment, error)
	GetRemoteOrgAssignmentByID(ctx context.Context, id string) (*core.RemoteOrgAssignment, error)
	GetRemoteOrgAssignmentsByOrganization(ctx context.Context, orgID string, limit, offset int) ([]*core.RemoteOrgAssignment, error)
	GetRemoteOrgAssignmentsByRemoteOrg(ctx context.Context, remoteOrgID string) ([]*core.RemoteOrgAssignment, error)
	UpdateRemoteOrgAssignment(ctx context.Context, id string, proxyOrgID string, active bool, notes *string) (*core.RemoteOrgAssignment, error)
	CreateAssignment(ctx context.Context, remoteOrgID, proxyOrgID string, active bool, notes string) (*core.RemoteOrgAssignment, error)
}

type remoteRepo struct {
	pool *pgxpool.Pool
}

func NewRemoteRepo(pool *pgxpool.Pool) RemoteRepo {
	return &remoteRepo{
		pool: pool,
	}
}

func (r *remoteRepo) CreateRemoteOrgAssignment(ctx context.Context, remoteOrgID, proxyOrgID string, notes *string) (*core.RemoteOrgAssignment, error) {
	id := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO remote_org_assignments (id, remote_org_id, proxy_org_id, active, notes, created_at)
		VALUES ($1, $2, $3, true, $4, $5)
		RETURNING id, remote_org_id, proxy_org_id, active, notes, created_at
	`

	var assignment core.RemoteOrgAssignment
	err := r.pool.QueryRow(ctx, query, id, remoteOrgID, proxyOrgID, notes, now).Scan(
		&assignment.ID, &assignment.RemoteOrgID, &assignment.ProxyOrgID, &assignment.Active, &assignment.Notes, &assignment.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &assignment, nil
}

func (r *remoteRepo) GetRemoteOrgAssignmentByID(ctx context.Context, id string) (*core.RemoteOrgAssignment, error) {
	query := `
		SELECT id, remote_org_id, proxy_org_id, active, notes, created_at
		FROM remote_org_assignments
		WHERE id = $1
	`

	var assignment core.RemoteOrgAssignment
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&assignment.ID, &assignment.RemoteOrgID, &assignment.ProxyOrgID, &assignment.Active, &assignment.Notes, &assignment.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, core.NewNotFoundError(core.ErrMsgRemoteAssignmentNotFound)
		}
		return nil, err
	}

	return &assignment, nil
}

func (r *remoteRepo) GetRemoteOrgAssignmentsByOrganization(ctx context.Context, orgID string, limit, offset int) ([]*core.RemoteOrgAssignment, error) {
	query := `
		SELECT id, remote_org_id, proxy_org_id, active, notes, created_at
		FROM remote_org_assignments
		WHERE proxy_org_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assignments []*core.RemoteOrgAssignment
	for rows.Next() {
		var assignment core.RemoteOrgAssignment
		err := rows.Scan(
			&assignment.ID, &assignment.RemoteOrgID, &assignment.ProxyOrgID, &assignment.Active, &assignment.Notes, &assignment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, &assignment)
	}

	return assignments, nil
}

func (r *remoteRepo) GetRemoteOrgAssignmentsByRemoteOrg(ctx context.Context, remoteOrgID string) ([]*core.RemoteOrgAssignment, error) {
	query := `
		SELECT id, remote_org_id, proxy_org_id, active, notes, created_at
		FROM remote_org_assignments
		WHERE remote_org_id = $1 AND active = true
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, remoteOrgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assignments []*core.RemoteOrgAssignment
	for rows.Next() {
		var assignment core.RemoteOrgAssignment
		err := rows.Scan(
			&assignment.ID, &assignment.RemoteOrgID, &assignment.ProxyOrgID, &assignment.Active, &assignment.Notes, &assignment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, &assignment)
	}

	return assignments, nil
}

func (r *remoteRepo) UpdateRemoteOrgAssignment(ctx context.Context, id string, proxyOrgID string, active bool, notes *string) (*core.RemoteOrgAssignment, error) {
	query := `
		UPDATE remote_org_assignments 
		SET proxy_org_id = $1, active = $2, notes = $3
		WHERE id = $4
		RETURNING id, remote_org_id, proxy_org_id, active, notes, created_at
	`

	var assignment core.RemoteOrgAssignment
	err := r.pool.QueryRow(ctx, query, proxyOrgID, active, notes, id).Scan(
		&assignment.ID, &assignment.RemoteOrgID, &assignment.ProxyOrgID, &assignment.Active, &assignment.Notes, &assignment.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, core.NewNotFoundError(core.ErrMsgRemoteAssignmentNotFound)
		}
		return nil, err
	}

	return &assignment, nil
}

func (r *remoteRepo) CreateAssignment(ctx context.Context, remoteOrgID, proxyOrgID string, active bool, notes string) (*core.RemoteOrgAssignment, error) {
	return r.CreateRemoteOrgAssignment(ctx, remoteOrgID, proxyOrgID, &notes)
}
