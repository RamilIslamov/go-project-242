package code

import (
	"fmt"
	urfaveCli "github.com/urfave/cli/v2"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
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
		Action: func(c *urfaveCli.Context) error {
			if c.Args().Len() == 0 {
				return urfaveCli.Exit("please provide a path", 1)
			}
			path := c.Args().First()
			human := c.Bool("human")
			all := c.Bool("all")
			recursive := c.Bool("recursive")

			size, err := GetSize(path, recursive, all)
			if err != nil {
				return err
			}

			formatted := FormatSize(size, human)

			fmt.Printf("%s\t%s\n", formatted, path)

			return nil
		},
	}
}

func GetSize(path string, recursive, all bool) (int64, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return 0, err
	}

	if !all && isHidden(fi) {
		return 0, nil
	}

	if !fi.IsDir() {
		return fileSize(fi, all)
	}
	return dirSize(path, recursive, all)
}

func fileSize(fi os.FileInfo, all bool) (int64, error) {
	if !all {
		if strings.HasPrefix(fi.Name(), ".") {
			return 0, nil
		}
	}
	return fi.Size(), nil
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

func dirSize(path string, recursive, all bool) (int64, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return 0, err
	}

	var sum int64
	for _, e := range entries {
		full := filepath.Join(path, e.Name())

		info, err := e.Info()
		if err != nil {
			return 0, err
		}

		if !all && isHidden(info) {
			continue
		}

		mode := info.Mode()
		if mode&fs.ModeSymlink != 0 {
			continue
		}

		if info.IsDir() {
			if recursive {
				sz, err := dirSize(full, recursive, all)
				if err != nil {
					return 0, err
				}
				sum += sz
			}
			continue
		}

		sum += info.Size()
	}
	return sum, nil
}

func isHidden(fi os.FileInfo) bool {
	name := fi.Name()

	if runtime.GOOS != "windows" {
		return strings.HasPrefix(name, ".")
	}

	if d, ok := fi.Sys().(*syscall.Win32FileAttributeData); ok {
		if d.FileAttributes&syscall.FILE_ATTRIBUTE_HIDDEN != 0 {
			return true
		}
	}
	return strings.HasPrefix(name, ".")
}
