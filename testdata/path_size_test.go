package code

import (
	"code"
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

func TestGetSize_Hidden_File_True(t *testing.T) {
	dir := t.TempDir()
	file := writeFile(t, dir, ".f.txt", "hello")

	got, err := code.GetSize(file, true)
	if err != nil {
		t.Fatalf("GetSize(file) error = %v", err)
	}

	want := int64(5)
	if got != want {
		t.Fatalf("GetSize(file) = %q, want %q", got, want)
	}
}

func TestGetSize_Hidden_File_False(t *testing.T) {
	dir := t.TempDir()
	file := writeFile(t, dir, ".hidden.txt", "hello")

	got, err := code.GetSize(file, false)
	if err != nil {
		t.Fatalf("GetSize(file) error = %v", err)
	}

	want := int64(0)
	if got != want {
		t.Fatalf("GetSize(file) = %q, want %q", got, want)
	}
}

func TestGetSize_DirFirstLevelOnly_No_Hidden_Files(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "a.bin", "aaa")
	writeFile(t, dir, "b.bin", "bbbb")
	writeFile(t, dir, "c.bin", "cc")
	sub := filepath.Join(dir, "sub")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatalf("Mkdir(sub): %v", err)
	}
	writeFile(t, sub, "deep.txt", "xxxxxxxxxx")

	got, err := code.GetSize(dir, false)
	if err != nil {
		t.Fatalf("GetSize(dir) error = %v", err)
	}

	want := int64(9)
	if got != want {
		t.Fatalf("GetSize(dir) = %q, want %q", got, want)
	}
}

func TestGetSize_DirFirstLevelOnly_With_Hidden_Files(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "a.bin", "aaa")
	writeFile(t, dir, ".b.bin", "bbbb")
	writeFile(t, dir, ".c.bin", "cc")
	sub := filepath.Join(dir, "sub")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatalf("Mkdir(sub): %v", err)
	}
	writeFile(t, sub, "deep.txt", "xxxxxxxxxx")

	got, err := code.GetSize(dir, false)
	if err != nil {
		t.Fatalf("GetSize(dir) error = %v", err)
	}

	want := int64(3)
	if got != want {
		t.Fatalf("GetSize(dir) = %q, want %q", got, want)
	}
}

func TestGetSize_PathNotExist(t *testing.T) {
	_, err := code.GetSize(filepath.Join(t.TempDir(), "nope.txt"), false)
	if err == nil {
		t.Fatalf("expected error for non-existent path, got nil")
	}
}

func TestFormatSize(t *testing.T) {
	const (
		KB = int64(1024)
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)

	tests := []struct {
		name  string
		size  int64
		human bool
		want  string
	}{
		{"bytes_non_human", 123, false, "123B"},
		{"zero_non_human", 0, false, "0B"},
		{"zero_human", 0, true, "0B"},

		{"999B_human", 999, true, "999B"},
		{"1023B_human", 1023, true, "1023B"},

		{"1KB", 1 * KB, true, "1KB"},
		{"2KB", 2 * KB, true, "2KB"},
		{"1MB", 1 * MB, true, "1MB"},

		{"1.5KB", 1536, true, "1.5KB"},
		{"round_to_2_0KB", 1997, true, "2.0KB"},

		{"10MB", 10 * MB, true, "10MB"},
		{"5GB", 5 * GB, true, "5GB"},
		{"5TB", 5 * TB, true, "5TB"},

		{"large_non_human", 1536, false, "1536B"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := code.FormatSize(tt.size, tt.human)
			if got != tt.want {
				t.Fatalf("FormatSize(%d, %v) = %q, want %q", tt.size, tt.human, got, tt.want)
			}
		})
	}
}
