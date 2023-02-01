package app

import (
	"fmt"

	"github.com/RipperAcskt/innotaxi/config"
	"github.com/RipperAcskt/innotaxi/internal/handler"
	"github.com/RipperAcskt/innotaxi/internal/repo/postgres"
	"github.com/RipperAcskt/innotaxi/internal/repo/redis"
	"github.com/RipperAcskt/innotaxi/internal/server"
	"github.com/RipperAcskt/innotaxi/internal/service"
	"github.com/sirupsen/logrus"

	"github.com/golang-migrate/migrate/v4"
)

func Run(log *logrus.Logger, cfg *config.Config) error {

	postgres, err := postgres.New(cfg)
	if err != nil {
		return fmt.Errorf("postgres new failed: %w", err)
	}
	defer postgres.Close()

	err = postgres.Migrate.Up()
	if err != migrate.ErrNoChange && err != nil {
		return fmt.Errorf("migrate up failed: %w", err)
	}

	redis, err := redis.New(cfg)
	if err != nil {
		return fmt.Errorf("redis new failed: %w", err)
	}
	defer redis.Close()

	service := service.New(postgres, redis, cfg.SALT, cfg)
	handler := handler.New(service, cfg, log)
	server := new(server.Server)
	if err := server.Run(handler.InitRouters(), cfg); err != nil {
		return fmt.Errorf("server run failed: %w", err)
	}
	return nil
}
