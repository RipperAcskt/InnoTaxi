package main

import (
	"github.com/RipperAcskt/innotaxi/internal/app"

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

	if err := app.Run(log); err != nil {
		log.Fatalf("app run failed: %v", err)
	}
}
