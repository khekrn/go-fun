package domain

import "context"

type AccountRepository interface {
	Create(ctx context.Context, account *Account) error
	Activate(ctx context.Context, id uint) error
	CreateAPIKey(ctx context.Context, accountId uint, name string) (string, error)
	DeactivateAPIKey(ctx context.Context, accountId uint, apiKey string) error
}

type ShortURLRepository interface {
	CreateURL(ctx context.Context, accountId, apiKeyId uint, sourceURL, shortCode string) (*ShortUrl, error)
	GetSourceURL(ctx context.Context, code string) (*ShortUrl, error)
	IncrementClicks(ctx context.Context, id uint) error
	DeactivateURL(ctx context.Context, accountId uint, code string) error
}
