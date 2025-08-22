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
	}
}

func GetSize(path string) (string, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return "", err
	}
	if !fi.IsDir() {
		return fmt.Sprintf("%dB\t%s\n", fi.Size(), path), nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}

	var sum int64
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		info, err := e.Info()
		if err != nil {
			return "", err
		}
		sum += info.Size()
	}
	return fmt.Sprintf("%dB\t%s\n", sum, path), nil
}
