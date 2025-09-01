package code

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func GetPathSize(path string, recursive, human, all bool) (string, error) {
	size, err := getSize(path, recursive, all)
	if err != nil {
		return "", err
	}
	return formatSize(size, human), nil
}

func getSize(path string, recursive, all bool) (int64, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return 0, err
	}
	if !all && isHiddenInfo(fi) {
		return 0, nil
	}
	if !fi.IsDir() {
		return fileSize(fi, all)
	}
	return dirSize(path, recursive, all)
}

func fileSize(fi os.FileInfo, all bool) (int64, error) {
	if !all {
		if isHiddenInfo(fi) {
			return 0, nil
		}
	}
	return fi.Size(), nil
}

func formatSize(size int64, human bool) string {
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

	if i == 0 {
		return fmt.Sprintf("%d%s", int64(val), units[i])
	}
	return fmt.Sprintf("%.1f%s", val, units[i])
}

func dirSize(path string, recursive, all bool) (int64, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		warnf("cannot read dir %q: %v", path, err)
		return 0, nil
	}

	var sum int64
	for _, e := range entries {
		info, subErr := e.Info()
		if subErr != nil {
			warnf("cannot stat %q: %v", filepath.Join(path, e.Name()), subErr)
			continue
		}
		if !all && isHiddenInfo(info) {
			continue
		}

		full := filepath.Join(path, info.Name())

		if info.Mode()&fs.ModeSymlink != 0 {
			continue
		}

		if info.IsDir() {
			if recursive {
				sz, dirErr := dirSize(full, recursive, all)
				if dirErr != nil {
					warnf("cannot descend into %q: %v", full, dirErr)
					continue
				}
				sum += sz
			}
			continue
		}

		sum += info.Size()
	}
	return sum, nil
}

func isHiddenInfo(fi os.FileInfo) bool {
	return strings.HasPrefix(fi.Name(), ".")
}

func warnf(format string, args ...any) {
	_, err := fmt.Fprintf(os.Stderr, "warning: "+format+"\n", args...)
	if err != nil {
		return
	}
}
