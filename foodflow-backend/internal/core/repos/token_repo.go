package repos

import (
	"context"
	"time"

	"foodflow/internal/core"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenRepo interface {
	CreateToken(ctx context.Context, collabID string, month string, tokensEarned int) (*core.Token, error)
	GetTokensByCollaborator(ctx context.Context, collabID string) ([]*core.Token, error)
	GetTokenHistory(ctx context.Context, collabID string, limit, offset int) ([]*core.Token, error)
	GetTotalTokens(ctx context.Context, collabID string) (int, error)
	GetTokenByCollabAndMonth(ctx context.Context, collabID string, month string) (*core.Token, error)
	UpdateTokenEarned(ctx context.Context, tokenID string, tokensEarned int) error
	RedeemTokens(ctx context.Context, collabID, month string, amount int, reason string) (*core.Token, error)
}

type tokenRepo struct {
	pool *pgxpool.Pool
}

func NewTokenRepo(pool *pgxpool.Pool) TokenRepo {
	return &tokenRepo{
		pool: pool,
	}
}

func (r *tokenRepo) CreateToken(ctx context.Context, collabID string, month string, tokensEarned int) (*core.Token, error) {
	id := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO tokens (id, collab_id, month, tokens_earned, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, collab_id, month, tokens_earned, created_at
	`

	var token core.Token
	err := r.pool.QueryRow(ctx, query, id, collabID, month, tokensEarned, now).Scan(
		&token.ID, &token.CollabID, &token.Month, &token.TokensEarned, &token.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *tokenRepo) GetTokensByCollaborator(ctx context.Context, collabID string) ([]*core.Token, error) {
	query := `
		SELECT id, collab_id, month, tokens_earned, created_at
		FROM tokens
		WHERE collab_id = $1 AND tokens_earned > 0
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, collabID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []*core.Token
	for rows.Next() {
		var token core.Token
		err := rows.Scan(
			&token.ID, &token.CollabID, &token.Month, &token.TokensEarned, &token.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, &token)
	}

	return tokens, nil
}

func (r *tokenRepo) GetTokenHistory(ctx context.Context, collabID string, limit, offset int) ([]*core.Token, error) {
	query := `
		SELECT id, collab_id, month, tokens_earned, created_at
		FROM tokens
		WHERE collab_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, collabID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []*core.Token
	for rows.Next() {
		var token core.Token
		err := rows.Scan(
			&token.ID, &token.CollabID, &token.Month, &token.TokensEarned, &token.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, &token)
	}

	return tokens, nil
}

func (r *tokenRepo) GetTotalTokens(ctx context.Context, collabID string) (int, error) {
	query := `
		SELECT COALESCE(SUM(tokens_earned), 0) as total
		FROM tokens
		WHERE collab_id = $1
	`

	var total int
	err := r.pool.QueryRow(ctx, query, collabID).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *tokenRepo) GetTokenByCollabAndMonth(ctx context.Context, collabID string, month string) (*core.Token, error) {
	query := `
		SELECT id, collab_id, month, tokens_earned, created_at
		FROM tokens
		WHERE collab_id = $1 AND month = $2
		ORDER BY created_at DESC
		LIMIT 1
	`

	var token core.Token
	err := r.pool.QueryRow(ctx, query, collabID, month).Scan(
		&token.ID, &token.CollabID, &token.Month, &token.TokensEarned, &token.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, core.NewNotFoundError(core.ErrMsgTokenNotFound)
		}
		return nil, err
	}

	return &token, nil
}

func (r *tokenRepo) UpdateTokenEarned(ctx context.Context, tokenID string, tokensEarned int) error {
	query := `
		UPDATE tokens 
		SET tokens_earned = $1
		WHERE id = $2
	`

	_, err := r.pool.Exec(ctx, query, tokensEarned, tokenID)
	if err != nil {
		return err
	}

	return nil
}

func (r *tokenRepo) RedeemTokens(ctx context.Context, collabID, month string, amount int, reason string) (*core.Token, error) {
	// For simplicity, create a negative token entry to represent redemption
	// In a real implementation, this might be a separate transaction table
	id := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO tokens (id, collab_id, month, tokens_earned, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, collab_id, month, tokens_earned, created_at
	`

	var token core.Token
	err := r.pool.QueryRow(ctx, query, id, collabID, month, -amount, now).Scan(
		&token.ID, &token.CollabID, &token.Month, &token.TokensEarned, &token.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
