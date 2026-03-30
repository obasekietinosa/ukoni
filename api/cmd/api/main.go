package main

import (
	"log/slog"
	"os"

	"ukoni/internal/config"
	"ukoni/internal/database"
	"ukoni/internal/mailer"
	"ukoni/internal/server"
)

// @title Ukoni API
// @version 1.0
// @description API for Ukoni application.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@ukoni.app

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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

	m := mailer.NewSMTPMailer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPFrom)

	srv := server.New(cfg, dbService, logger, m)
	if err := srv.Run(); err != nil {
		logger.Error("server failed to run", "error", err)
		os.Exit(1)
	}
}
