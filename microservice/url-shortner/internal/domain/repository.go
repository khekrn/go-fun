package domain

type AccountRepository interface {
	Create(account *Account) error
	Activate(id uint) error
	CreateAPIKey(accountId uint, name string) (string, error)
	DeactivateAPIKey(accountId uint, apiKey string) error
}

type ShortURLRepository interface {
	CreateURL(accountId, apiKeyId uint, sourceURL string, customSlug string) (*ShortUrl, error)
	GetSourceURL(code string) (*ShortUrl, error)
	IncrementClicks(id uint) error
	DeactivateURL(accountId uint, code string) error
}
