package main

import (
	"fmt"
	"log"

	"github.com/OtavioPompolini/project-postman/app"
)

func main() {
	fmt.Println("PUDIM")

	app, err := app.NewApp()
	if err != nil {
		log.Panicln(err)
	}

	if err := app.Run(); err != nil {
		log.Panicln(err)
	}
}
