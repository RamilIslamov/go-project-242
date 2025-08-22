package code

import (
	"code"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeFile: %v", err)
	}
	return p
}

func TestGetSize_File(t *testing.T) {
	dir := t.TempDir()
	file := writeFile(t, dir, "f.txt", "hello")

	got, err := code.GetSize(file)
	if err != nil {
		t.Fatalf("GetSize(file) error = %v", err)
	}

	want := fmt.Sprintf("%dB\t%s\n", int64(5), file)
	if got != want {
		t.Fatalf("GetSize(file) = %q, want %q", got, want)
	}
}

func TestGetSize_DirFirstLevelOnly(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "a.bin", "aaa")
	writeFile(t, dir, "b.bin", "bbbb")
	sub := filepath.Join(dir, "sub")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatalf("Mkdir(sub): %v", err)
	}
	writeFile(t, sub, "deep.txt", "xxxxxxxxxx")

	got, err := code.GetSize(dir)
	if err != nil {
		t.Fatalf("GetSize(dir) error = %v", err)
	}

	want := fmt.Sprintf("%dB\t%s\n", int64(7), dir)
	if got != want {
		t.Fatalf("GetSize(dir) = %q, want %q", got, want)
	}
}

func TestGetSize_PathNotExist(t *testing.T) {
	_, err := code.GetSize(filepath.Join(t.TempDir(), "nope.txt"))
	if err == nil {
		t.Fatalf("expected error for non-existent path, got nil")
	}
}
