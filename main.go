package main

import (
	"log"
)

func main() {

	app, err := NewApp();
	if err != nil {
		log.Panicln(err)
	}

	if err := app.Run(); err != nil {
		log.Panicln(err)
	}
}
