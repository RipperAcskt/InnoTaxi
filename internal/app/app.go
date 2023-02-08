package app

import (
	"fmt"
	"os"

	"github.com/RipperAcskt/innotaxi/config"
	"github.com/RipperAcskt/innotaxi/internal/handler"
	"github.com/RipperAcskt/innotaxi/internal/repo/mongo"
	"github.com/RipperAcskt/innotaxi/internal/repo/postgres"
	"github.com/RipperAcskt/innotaxi/internal/repo/redis"
	"github.com/RipperAcskt/innotaxi/internal/server"
	"github.com/RipperAcskt/innotaxi/internal/service"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/golang-migrate/migrate/v4"
)

func Run() error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("config new failed: %w", err)
	}

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

	mongo, err := mongo.New(cfg)
	if err != nil {
		return fmt.Errorf("mongo new failed: %w", err)
	}
	defer mongo.Close()

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	writer := zapcore.AddSync(mongo)
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)
	log := zap.New(core, zap.AddCaller())
	defer log.Sync()

	service := service.New(postgres, redis, cfg.SALT, cfg)
	handler := handler.New(service, cfg, log)
	server := new(server.Server)
	if err := server.Run(handler.InitRouters(), cfg); err != nil {
		return fmt.Errorf("server run failed: %w", err)
	}
	return nil
}
