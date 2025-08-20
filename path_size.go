package pathsize

import urfaveCli "github.com/urfave/cli/v2"

func NewApp() *urfaveCli.App {
	return &urfaveCli.App{
		Name:  "hexlet-path-size",
		Usage: "print size of a file or directory",
	}
}
