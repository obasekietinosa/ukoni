package main

import (
	"log/slog"
	"os"

	"ukoni/internal/config"
	"ukoni/internal/database"
	"ukoni/internal/server"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	dbService, err := database.New(cfg.DBURL)
	if err != nil {
		logger.Error("failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer dbService.Close()
	logger.Info("database connected")

	srv := server.New(cfg, dbService, logger)
	if err := srv.Run(); err != nil {
		logger.Error("server failed to run", "error", err)
		os.Exit(1)
	}
}
