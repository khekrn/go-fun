package main

import (
	"coding2fun.in/url-shortner/database"
	"coding2fun.in/url-shortner/internal/config"
	"coding2fun.in/url-shortner/internal/log"
	"coding2fun.in/url-shortner/internal/server"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load("resources/local.ini")
	if err != nil {
		panic("Failed to load config " + err.Error())
	}

	log.InitLogger(cfg)
	defer log.Sync()

	dbService, err := database.NewService(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database: ", zap.Error(err))
	}
	defer dbService.Close()

	db := dbService.GetConnection()
	err = dbService.Migrate()
	if err != nil {
		log.Fatal("Failed to run database migrations", zap.Error(err))
	}

	srv := server.NewServer(cfg, db)

	// Start server with graceful shutdown handling
	if err := srv.StartWithGracefulShutdown(); err != nil {
		log.Fatal("Server error", zap.Error(err))
	}

	log.Info("Server exited properly")

}
