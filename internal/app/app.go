package app

import (
	"fmt"
	"os"

	"github.com/RipperAcskt/innotaxi/internal/handler"
	"github.com/RipperAcskt/innotaxi/internal/repo/postgres"
	"github.com/RipperAcskt/innotaxi/internal/server"
	"github.com/RipperAcskt/innotaxi/internal/service"
	"github.com/golang-migrate/migrate/v4"
)

func Run() error {
	urlDB := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("DBUSERNAME"), os.Getenv("DBPASSWORD"), os.Getenv("DBHOST"), os.Getenv("DBPORT"), os.Getenv("DBNAME"))
	db, err := postgres.New(urlDB)
	if err != nil {
		return fmt.Errorf("postgres new failed: %w", err)
	}
	defer db.Close()

	err = db.Migrate.Up()
	if err != migrate.ErrNoChange && err != nil {
		return fmt.Errorf("migrate up failed: %w", err)
	}

	service := service.New(db, os.Getenv("SALT"))
	handler := handler.New(service)
	server := new(server.Server)
	if err := server.Run(handler.InitRouters()); err != nil {
		return fmt.Errorf("server run failed: %w", err)
	}
	return nil
}
