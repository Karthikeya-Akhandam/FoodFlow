package repos

import (
	"context"
	"time"

	"foodflow/internal/core"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CollabRepo interface {
	CreateCollaborator(ctx context.Context, userID string, collabType string) (*core.Collaborator, error)
	GetCollaboratorByID(ctx context.Context, id string) (*core.Collaborator, error)
	GetCollaboratorByUserID(ctx context.Context, userID string) (*core.Collaborator, error)
	UpdateCollaborator(ctx context.Context, id string, collabType string) (*core.Collaborator, error)
	ListCollaborators(ctx context.Context, limit, offset int) ([]*core.Collaborator, error)
	GetAllCollaborators(ctx context.Context, limit, offset int) ([]*core.Collaborator, error)
}

type collabRepo struct {
	pool *pgxpool.Pool
}

func NewCollabRepo(pool *pgxpool.Pool) CollabRepo {
	return &collabRepo{
		pool: pool,
	}
}

func (r *collabRepo) CreateCollaborator(ctx context.Context, userID string, collabType string) (*core.Collaborator, error) {
	id := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO collaborators (id, user_id, collab_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, collab_type, created_at, updated_at
	`

	var collab core.Collaborator
	err := r.pool.QueryRow(ctx, query, id, userID, collabType, now, now).Scan(
		&collab.ID, &collab.UserID, &collab.CollabType, &collab.CreatedAt, &collab.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &collab, nil
}

func (r *collabRepo) GetCollaboratorByID(ctx context.Context, id string) (*core.Collaborator, error) {
	query := `
		SELECT id, user_id, collab_type, created_at, updated_at
		FROM collaborators
		WHERE id = $1
	`

	var collab core.Collaborator
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&collab.ID, &collab.UserID, &collab.CollabType, &collab.CreatedAt, &collab.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, core.NewNotFoundError(core.ErrMsgCollaboratorNotFound)
		}
		return nil, err
	}

	return &collab, nil
}

func (r *collabRepo) GetCollaboratorByUserID(ctx context.Context, userID string) (*core.Collaborator, error) {
	query := `
		SELECT id, user_id, collab_type, created_at, updated_at
		FROM collaborators
		WHERE user_id = $1
	`

	var collab core.Collaborator
	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&collab.ID, &collab.UserID, &collab.CollabType, &collab.CreatedAt, &collab.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, core.NewNotFoundError(core.ErrMsgCollaboratorNotFound)
		}
		return nil, err
	}

	return &collab, nil
}

func (r *collabRepo) UpdateCollaborator(ctx context.Context, id string, collabType string) (*core.Collaborator, error) {
	query := `
		UPDATE collaborators 
		SET collab_type = $1, updated_at = $2
		WHERE id = $3
		RETURNING id, user_id, collab_type, created_at, updated_at
	`

	var collab core.Collaborator
	err := r.pool.QueryRow(ctx, query, collabType, time.Now(), id).Scan(
		&collab.ID, &collab.UserID, &collab.CollabType, &collab.CreatedAt, &collab.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, core.NewNotFoundError(core.ErrMsgCollaboratorNotFound)
		}
		return nil, err
	}

	return &collab, nil
}

func (r *collabRepo) ListCollaborators(ctx context.Context, limit, offset int) ([]*core.Collaborator, error) {
	query := `
		SELECT id, user_id, collab_type, created_at, updated_at
		FROM collaborators
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collabs []*core.Collaborator
	for rows.Next() {
		var collab core.Collaborator
		err := rows.Scan(
			&collab.ID, &collab.UserID, &collab.CollabType, &collab.CreatedAt, &collab.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		collabs = append(collabs, &collab)
	}

	return collabs, nil
}

func (r *collabRepo) GetAllCollaborators(ctx context.Context, limit, offset int) ([]*core.Collaborator, error) {
	return r.ListCollaborators(ctx, limit, offset)
}
