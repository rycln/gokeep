package main

import (
	"log"

	"github.com/rycln/gokeep/internal/server/app"
)

func main() {
	app, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
