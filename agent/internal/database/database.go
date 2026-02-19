package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Service interface {
	Health() map[string]string
	Close() error
	GetDB() *sql.DB
}

type service struct {
	db *sql.DB
}

func New(connStr string) (Service, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &service{db: db}, nil
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		slog.Error("db down", "error", err)
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "It's healthy"
	stats["open_connections"] = fmt.Sprintf("%d", s.db.Stats().OpenConnections)
	stats["in_use"] = fmt.Sprintf("%d", s.db.Stats().InUse)
	stats["idle"] = fmt.Sprintf("%d", s.db.Stats().Idle)
	stats["wait_count"] = fmt.Sprintf("%d", s.db.Stats().WaitCount)
	stats["wait_duration"] = fmt.Sprintf("%v", s.db.Stats().WaitDuration)
	stats["max_idle_closed"] = fmt.Sprintf("%d", s.db.Stats().MaxIdleClosed)
	stats["max_lifetime_closed"] = fmt.Sprintf("%d", s.db.Stats().MaxLifetimeClosed)

	return stats
}

func (s *service) Close() error {
	return s.db.Close()
}

func (s *service) GetDB() *sql.DB {
	return s.db
}
