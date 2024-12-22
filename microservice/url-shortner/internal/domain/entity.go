package domain

import (
	"gorm.io/gorm"
	"time"
)

type Account struct {
	gorm.Model
	Email     string     `gorm:"uniqueIndex;not null"`
	IsActive  bool       `gorm:"default:false"`
	APIKeys   []APIKey   `gorm:"foreignKey:AccountId"`
	ShortUrls []ShortUrl `gorm:"foreignKey:AccountId"`
}

type APIKey struct {
	gorm.Model
	AccountId uint
	Key       string `gorm:"uniqueIndex;not null"`
	Name      string
	IsActive  bool `gorm:"default:true"`
	LastUsed  time.Time
	ExpiresAt time.Time
}

type ShortUrl struct {
	gorm.Model
	AccountId     uint
	APIKeyId      uint
	OriginalURL   string `gorm:"not null"`
	ShortCode     string `gorm:"uniqueIndex;not null"`
	ExpiresAt     time.Time
	IsActive      bool  `gorm:"default:true"`
	Clicks        int64 `gorm:"default:0"`
	LastClickedAt time.Time
	CustomSlug    string `gorm:"uniqueIndex"`
}

type URLAnalytics struct {
	gorm.Model
	ShortURLId    uint  `gorm:"uniqueIndex"`
	TotalClicks   int64 `gorm:"default:0"`
	LastClickedAt time.Time
}
