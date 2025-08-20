package main

import (
	"log"
	"os"

	pathsize "code"
)

func main() {
	app := pathsize.NewApp()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
