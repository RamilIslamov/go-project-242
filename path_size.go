package code

import (
	"fmt"
	urfaveCli "github.com/urfave/cli/v2"
	"os"
)

func NewApp() *urfaveCli.App {
	return &urfaveCli.App{
		Name:  "hexlet-path-size",
		Usage: "print size of a file or directory",
		Flags: []urfaveCli.Flag{
			&urfaveCli.BoolFlag{
				Name:    "human",
				Aliases: []string{"H"},
				Usage:   "human readable sizes",
			},
		},
		Action: func(c *urfaveCli.Context) error {
			if c.Args().Len() == 0 {
				return urfaveCli.Exit("please provide a path", 1)
			}
			path := c.Args().First()
			human := c.Bool("human")

			size, err := GetSize(path)
			if err != nil {
				return err
			}

			formatted := FormatSize(size, human)

			fmt.Printf("%s\t%s\n", formatted, path)

			return nil
		},
	}
}

func GetSize(path string) (int64, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return 0, err
	}
	if !fi.IsDir() {
		return fi.Size(), nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return 0, err
	}

	var sum int64
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		info, err := e.Info()
		if err != nil {
			return 0, err
		}
		sum += info.Size()
	}
	return sum, nil
}

func FormatSize(size int64, human bool) string {
	if !human {
		return fmt.Sprintf("%dB", size)
	}
	const step = 1024.0
	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}

	val := float64(size)
	i := 0
	for val >= step && i < len(units)-1 {
		val /= step
		i++
	}

	if val == float64(int64(val)) {
		return fmt.Sprintf("%.0f%s", val, units[i])
	}
	return fmt.Sprintf("%.1f%s", val, units[i])
}
