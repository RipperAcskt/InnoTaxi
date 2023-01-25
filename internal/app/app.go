package app

import (
	"fmt"

	"github.com/RipperAcskt/innotaxi/config"
	"github.com/RipperAcskt/innotaxi/internal/handler"
	"github.com/RipperAcskt/innotaxi/internal/repo/postgres"
	"github.com/RipperAcskt/innotaxi/internal/server"
	"github.com/RipperAcskt/innotaxi/internal/service"

	"github.com/golang-migrate/migrate/v4"
)

func Run() error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("config new failed: %w", err)
	}

	db, err := postgres.New(cfg.GetDBUrl(), cfg)
	if err != nil {
		return fmt.Errorf("postgres new failed: %w", err)
	}
	defer db.Close()

	err = db.Migrate.Up()
	if err != migrate.ErrNoChange && err != nil {
		return fmt.Errorf("migrate up failed: %w", err)
	}

	service := service.New(db, cfg.SALT, cfg)
	handler := handler.New(service, cfg)
	server := new(server.Server)
	if err := server.Run(handler.InitRouters(), cfg); err != nil {
		return fmt.Errorf("server run failed: %w", err)
	}
	return nil
}
