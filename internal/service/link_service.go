package service

import (
	"context"
	"errors"
	"pht/pet/link_shortener/internal/domain"
	"pht/pet/link_shortener/internal/domain/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LinkService struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewLinkService(pool *pgxpool.Pool, queries *db.Queries) *LinkService {
	return &LinkService{pool: pool, queries: queries}
}

func (ls *LinkService) Save(ctx context.Context, url string, code string) (*domain.Link, error) {
	createLinkParams := db.CreateLinkParams{ShortCode: code, LongUrl: url}

	savedLink, err := ls.queries.CreateLink(ctx, createLinkParams)
	if err != nil {
		return nil, errors.New("cannot save the link")
	}

	link := &domain.Link{
		ID:        savedLink.ID,
		ShortCode: savedLink.ShortCode,
		LongUrl:   savedLink.LongUrl,
		Clicks:    savedLink.Clicks,
	}

	return link, nil
}

func (ls *LinkService) GetLink(ctx context.Context, code string) (*domain.Link, error) {
	linkFromDB, err := ls.queries.GetLinkByCode(ctx, code)
	if err != nil {
		return nil, errors.New("not found")
	}

	link := &domain.Link{
		ID:        linkFromDB.ID,
		ShortCode: linkFromDB.ShortCode,
		LongUrl:   linkFromDB.LongUrl,
		Clicks:    linkFromDB.Clicks,
	}

	return link, nil
}
