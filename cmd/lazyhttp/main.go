package main

import (
	"log"

	"github.com/OtavioPompolini/project-postman/internal/app"
)

func main() {
	app, err := app.NewApp()
	if err != nil {
		log.Panicln(err)
	}

	if err := app.Run(); err != nil {
		log.Panicln(err)
	}
}
