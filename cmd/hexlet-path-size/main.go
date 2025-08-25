package main

import (
	"code"
	"os"
)

func main() {
	app := code.NewApp()
	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
