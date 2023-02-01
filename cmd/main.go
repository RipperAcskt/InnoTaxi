package main

import (
	"github.com/RipperAcskt/innotaxi/config"
	"github.com/RipperAcskt/innotaxi/internal/app"
	"github.com/RipperAcskt/innotaxi/internal/repo/mongo"

	"github.com/sirupsen/logrus"
)

// @title InnoTaxi API
// @version 1.0
// @description API for order taxi
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  ripper@gmail.com

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	log := logrus.New()

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config new failed: %w", err)
	}

	mongo, err := mongo.New(cfg)
	if err != nil {
		log.Fatalf("mongo new failed: %w", err)
	}
	log.SetOutput(mongo)

	if err := app.Run(log, cfg); err != nil {
		log.Fatalf("app run failed: %v", err)
	}
}
