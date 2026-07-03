package service

import (
	"context"
	"fmt"
	"pht/pet/link_shortener/internal/domain"
	"pht/pet/link_shortener/internal/domain/db"
)

type URLRepository interface {
	SaveLink(context.Context, db.CreateLinkParams) (*domain.Link, error)
	GetURLAndIncrementLinkClicks(context.Context, string) (string, error)
	GetClicks(context.Context, string) (int32, error)
}

type LinkService struct {
	querier URLRepository
}

func NewLinkService(querier URLRepository) *LinkService {
	return &LinkService{querier: querier}
}

func (ls *LinkService) Save(ctx context.Context, url, code string) (*domain.Link, error) {
	data := domain.Link{ShortCode: code, LongUrl: url}
	if err := data.Validate(); err != nil {
		return nil, err
	}

	createLinkParams := db.CreateLinkParams{ShortCode: data.ShortCode, LongUrl: data.LongUrl}
	savedLink, err := ls.querier.SaveLink(ctx, createLinkParams)
	if err != nil {
		return nil, fmt.Errorf("failed to save the link, err: %w", err)
	}

	link := &domain.Link{
		ID:        savedLink.ID,
		ShortCode: savedLink.ShortCode,
		LongUrl:   savedLink.LongUrl,
		Clicks:    savedLink.Clicks,
	}

	return link, nil
}

func (ls *LinkService) GetURL(ctx context.Context, code string) (string, error) {
	url, err := ls.querier.GetURLAndIncrementLinkClicks(ctx, code)
	if err != nil {
		return "", fmt.Errorf("failed to get the url: %w", err)
	}

	return url, nil
}

func (ls *LinkService) GetClicks(ctx context.Context, code string) (int32, error) {
	clicks, err := ls.querier.GetClicks(ctx, code)
	if err != nil {
		return 0, fmt.Errorf("failed to get clicks: %w", err)
	}

	return clicks, nil
}
