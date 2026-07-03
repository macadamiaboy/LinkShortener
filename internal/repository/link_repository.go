package repository

import (
	"context"
	"errors"
	"pht/pet/link_shortener/internal/domain"
	"pht/pet/link_shortener/internal/domain/db"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PGXURLRepository struct {
	q *db.Queries
}

func NewPGXURLRepository(q *db.Queries) *PGXURLRepository {
	return &PGXURLRepository{q: q}
}

func (repo *PGXURLRepository) SaveLink(ctx context.Context, params db.CreateLinkParams) (*domain.Link, error) {
	linkRow, err := repo.q.CreateLink(ctx, params)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, domain.ErrCodeUniquenessConflict
		}
		return nil, domain.ErrInternal
	}

	link := &domain.Link{ID: linkRow.ID, ShortCode: linkRow.ShortCode, LongUrl: linkRow.LongUrl, Clicks: linkRow.Clicks}
	return link, nil
}

func (repo *PGXURLRepository) GetURLAndIncrementLinkClicks(ctx context.Context, code string) (string, error) {
	url, err := repo.q.GetURLAndIncrementLinkClicks(ctx, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", domain.ErrLinkNotFound
		}
		return "", domain.ErrInternal
	}

	return url, nil
}

func (repo *PGXURLRepository) GetClicks(ctx context.Context, code string) (int32, error) {
	clicks, err := repo.q.GetClicksByCode(ctx, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, domain.ErrLinkNotFound
		}
		return 0, domain.ErrInternal
	}

	return clicks, nil
}
