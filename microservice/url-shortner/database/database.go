package database

import (
	"coding2fun.in/url-shortner/internal/config"
	"coding2fun.in/url-shortner/internal/log"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"time"
)

type Service interface {
	GetConnection() *gorm.DB
	Migrate() error
	Close() error
}

type service struct {
	db *gorm.DB
}

func NewService(cfg *config.DatabaseConfig) (Service, error) {
	log.Info("Connecting to database", zap.String("host", cfg.Host), zap.String("dbName", cfg.Name))
	db, err := gorm.Open(postgres.Open(cfg.ConnectionURL()), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		log.Error("Failed to connect to database", zap.Error(err))
		os.Exit(1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Error("Failed to get database instance", zap.Error(err))
		os.Exit(1)
	}

	sqlDB.SetConnMaxLifetime(1 * time.Hour)
	sqlDB.SetMaxOpenConns(64)

	return &service{
		db: db,
	}, nil
}

func (s *service) GetConnection() *gorm.DB {
	return s.db
}

func (s *service) Migrate() error {
	return nil
}

func (s *service) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}
	return sqlDB.Close()
}
