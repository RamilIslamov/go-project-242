package code

import (
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
	got, err := getSize("./testdata/.f.txt", false, true)
	if err != nil {
		t.Fatalf("getSize(file) error = %v", err)
	}

	want := int64(5)
	if got != want {
		t.Fatalf("getSize(file) = %d, want %d", got, want)
	}
}

func TestGetSize_Hidden_File_False(t *testing.T) {
	got, err := getSize("./testdata/.f.txt", false, false)
	if err != nil {
		t.Fatalf("getSize(file) error = %v", err)
	}

	want := int64(0)
	if got != want {
		t.Fatalf("getSize(file) = %d, want %d", got, want)
	}
}

func TestGetSize_DirFirstLevelOnly_No_Hidden_Files(t *testing.T) {
	got, err := getSize("./testdata/flevel", false, false)
	if err != nil {
		t.Fatalf("getSize(dir) error = %v", err)
	}

	want := int64(9)
	if got != want {
		t.Fatalf("getSize(dir) = %d, want %d", got, want)
	}
}

func TestGetSize_DirFirstLevelOnly_With_Hidden_Files(t *testing.T) {
	got, err := getSize("./testdata/flevel", false, false)
	if err != nil {
		t.Fatalf("getSize(dir) error = %v", err)
	}

	want := int64(9)
	if got != want {
		t.Fatalf("getSize(dir) = %d, want %d", got, want)
	}
}

func TestGetSize_DirAllLevels_With_Hidden_Files(t *testing.T) {
	got, err := getSize("./testdata/flevel", true, true)
	if err != nil {
		t.Fatalf("getSize(dir) error = %v", err)
	}

	want := int64(1465)
	if got != want {
		t.Fatalf("getSize(dir) = %d, want %d", got, want)
	}
}

func TestGetSize_PathNotExist(t *testing.T) {
	_, err := getSize("./testdata/flevel/nope.txt", false, false)
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

		{"1KB", 1 * KB, true, "1.0KB"},
		{"2KB", 2 * KB, true, "2.0KB"},
		{"1MB", 1 * MB, true, "1.0MB"},

		{"1.5KB", 1536, true, "1.5KB"},
		{"round_to_2_0KB", 1997, true, "2.0KB"},

		{"10MB", 10 * MB, true, "10.0MB"},
		{"5GB", 5 * GB, true, "5.0GB"},
		{"5TB", 5 * TB, true, "5.0TB"},

		{"large_non_human", 1536, false, "1536B"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatSize(tt.size, tt.human)
			if got != tt.want {
				t.Fatalf("FormatSize(%d, %v) = %q, want %q", tt.size, tt.human, got, tt.want)
			}
		})
	}
}

func TestGetPathSize_All_Flags(t *testing.T) {

	const path = "./testdata"

	tests := []struct {
		name      string
		recursive bool
		human     bool
		all       bool
		want      string
	}{
		{"true-true-true", true, true, true, "1.4KB"},
		{"true-true-false", true, true, false, "674B"},
		{"true-false-true", true, false, true, "1470B"},
		{"false-true-true", false, true, true, "5B"},
		{"false-false-false", false, false, false, "0B"},
		{"false-false-true", false, false, true, "5B"},
		{"false-true-false", false, true, false, "0B"},
		{"true-false-false", true, false, false, "674B"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPathSize(path, tt.recursive, tt.human, tt.all)
			if err != nil {
				t.Fatalf("GetPathSize(dir) error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("GetPathSize(dir) = %s, want %s", got, tt.want)
			}
		})
	}
}
