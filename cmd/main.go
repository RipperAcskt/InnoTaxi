package main

import (
	"log"

	"github.com/RipperAcskt/innotaxi/internal/handler"
	"github.com/RipperAcskt/innotaxi/internal/repo/postgres"
	"github.com/RipperAcskt/innotaxi/internal/server"
	"github.com/RipperAcskt/innotaxi/internal/service"
	"github.com/golang-migrate/migrate/v4"
)

func main() {
	db, err := postgres.New("postgres://ripper:150403@localhost:5432/innotaxi")
	if err != nil {
		log.Fatalf("postgres new failed: %v", err)
	}

	err = db.Migrate.Down()
	if err != migrate.ErrNoChange && err != nil {
		log.Fatalf("migrate down failed: %v", err)
	}

	err = db.Migrate.Up()
	if err != migrate.ErrNoChange && err != nil {
		log.Fatalf("migrate up failed: %v", err)
	}

	service := service.New(db)
	handler := handler.New(service)
	server := new(server.Server)
	if err := server.Run(handler.InitRouters()); err != nil {
		log.Fatalf("server run failed: %v", err)
	}
}
