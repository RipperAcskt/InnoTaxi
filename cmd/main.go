package main

import (
	"log"

	"github.com/RipperAcskt/innotaxi/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("app run failed: %v", err)
	}
}
