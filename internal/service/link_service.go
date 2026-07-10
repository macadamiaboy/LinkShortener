package service

import (
	"context"
	"errors"
	"fmt"
	"pht/pet/link_shortener/internal/domain"
	"pht/pet/link_shortener/internal/domain/db"
	"time"

	"github.com/redis/go-redis/v9"
)

type URLRepository interface {
	SaveLink(context.Context, db.CreateLinkParams) (*domain.Link, error)
	GetURLAndIncrementLinkClicks(context.Context, string) (string, error)
	GetClicks(context.Context, string) (int32, error)
}

type WarnLogger interface {
	Warn(msg string, args ...any)
}

type LinkService struct {
	querier URLRepository
	redis   *redis.Client
	logger  WarnLogger
}

func NewLinkService(querier URLRepository, redisClient *redis.Client, logger WarnLogger) *LinkService {
	return &LinkService{querier: querier, redis: redisClient, logger: logger}
}

func (ls *LinkService) Save(ctx context.Context, url, code string) (*domain.Link, error) {
	data := domain.Link{ShortCode: code, LongUrl: url}
	if err := data.Validate(); err != nil {
		return nil, err
	}

	createLinkParams := db.CreateLinkParams{ShortCode: data.ShortCode, LongUrl: data.LongUrl}
	savedLink, err := ls.querier.SaveLink(ctx, createLinkParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create the short link for url, err: %w", err)
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
	redisKey := fmt.Sprintf("url:%s", code)

	redisURL, err := ls.redis.Get(ctx, redisKey).Result()
	if err == nil {
		return redisURL, nil
	}

	if !errors.Is(err, redis.Nil) {
		ls.logger.Warn("redis error", "error", err)
	}

	url, err := ls.querier.GetURLAndIncrementLinkClicks(ctx, code)
	if err != nil {
		return "", fmt.Errorf("failed to get the url incrementing the clicks count: %w", err)
	}

	err = ls.redis.Set(ctx, redisKey, url, 24*time.Hour).Err()
	if err != nil {
		ls.logger.Warn("failed to cache the URL in Redis", "error", err)
	}

	return url, nil
}

func (ls *LinkService) GetClicks(ctx context.Context, code string) (int32, error) {
	clicks, err := ls.querier.GetClicks(ctx, code)
	if err != nil {
		return 0, fmt.Errorf("failed to get clicks count: %w", err)
	}

	return clicks, nil
}
