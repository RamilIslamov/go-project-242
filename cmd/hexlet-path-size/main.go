package main

import (
	"code"
	"context"
	"fmt"
	urfaveCli "github.com/urfave/cli/v3"
	"os"
)

func main() {
	app := newApp()

	if err := app.Run(context.Background(), os.Args); err != nil {
		os.Exit(1)
	}
}

func newApp() *urfaveCli.Command {
	return &urfaveCli.Command{
		Name:  "hexlet-path-size",
		Usage: "print size of a file or directory",
		Flags: []urfaveCli.Flag{
			&urfaveCli.BoolFlag{
				Name:    "human",
				Aliases: []string{"H"},
				Usage:   "human readable sizes",
			},
			&urfaveCli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "include hidden files and directories",
			},
			&urfaveCli.BoolFlag{
				Name:    "recursive",
				Aliases: []string{"r"},
				Usage:   "recursive size of directories",
			},
		},
		Action: func(_ context.Context, cmd *urfaveCli.Command) error {
			if cmd.Args().Len() == 0 {
				return urfaveCli.Exit("please provide a path", 1)
			}

			path := cmd.Args().First()
			recursive := cmd.Bool("recursive")
			all := cmd.Bool("all")
			human := cmd.Bool("human")

			size, err := code.GetPathSize(path, recursive, human, all)
			if err != nil {
				return err
			}

			fmt.Printf("%s\t%s\n", size, path)
			return nil
		},
	}
}
