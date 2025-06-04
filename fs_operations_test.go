package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestIsSafePath(t *testing.T) {
	baseDir := t.TempDir()

	safeSubDir := filepath.Join(baseDir, "safe-subdir")
	if err := os.Mkdir(safeSubDir, 0755); err != nil {
		t.Fatalf("Failed to create safe subdirectory: %v", err)
	}

	safeFile := filepath.Join(safeSubDir, "file.txt")
	if _, err := os.Create(safeFile); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	outsideDir := t.TempDir()

	tests := []struct {
		name     string
		base     string
		target   string
		expected bool
	}{
		{
			name:     "file inside base",
			base:     baseDir,
			target:   safeFile,
			expected: true,
		},
		{
			name:     "directory inside base",
			base:     baseDir,
			target:   safeSubDir,
			expected: true,
		},
		{
			name:     "file outside base",
			base:     baseDir,
			target:   filepath.Join(outsideDir, "file.txt"),
			expected: false,
		},
		{
			name:     "relative path within base",
			base:     baseDir,
			target:   "./safe-subdir/file.txt",
			expected: true,
		},
		{
			name:     "parent directory traversal",
			base:     safeSubDir,
			target:   filepath.Join(safeSubDir, "..", "safe-subdir", "file.txt"),
			expected: true,
		},
		{
			name:     "invalid path",
			base:     baseDir,
			target:   filepath.Join(baseDir, "nonexistent", "file.txt"),
			expected: true, // The path is still considered safe even if it doesn't exist
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For relative path tests, we need to change working directory
			if tt.name == "relative path within base" {
				oldWd, err := os.Getwd()
				if err != nil {
					t.Fatal(err)
				}
				defer os.Chdir(oldWd)
				if err := os.Chdir(baseDir); err != nil {
					t.Fatal(err)
				}
			}

			t.Logf("Testing for path '%s'", tt.target)

			actual := isSafePath(tt.base, tt.target)
			if actual != tt.expected {
				t.Errorf("isSafePath(%q, %q) = %v, want %v", tt.base, tt.target, actual, tt.expected)
			}
		})
	}
}

func TestAssertPath(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name     string
		path     string
		setup    func()
		wantInfo bool
		wantErr  bool
		wantBool bool
	}{
		{
			name:     "existing file",
			path:     filepath.Join(tmpDir, "exists.txt"),
			setup:    func() { os.WriteFile(filepath.Join(tmpDir, "exists.txt"), []byte("test"), 0644) },
			wantInfo: true,
			wantErr:  false,
			wantBool: true,
		},
		{
			name:     "non-existent file",
			path:     filepath.Join(tmpDir, "nonexistent.txt"),
			wantInfo: false,
			wantErr:  false,
			wantBool: false,
		},
		{
			name:     "existing directory",
			path:     tmpDir,
			wantInfo: true,
			wantErr:  false,
			wantBool: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			gotInfo, gotErr, gotBool := assertPath(tt.path)

			if tt.wantInfo && gotInfo == nil {
				t.Error("expected FileInfo, got nil")
			}
			if !tt.wantInfo && gotInfo != nil {
				t.Error("expected nil FileInfo, got value")
			}

			if tt.wantErr && gotErr == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && gotErr != nil {
				t.Errorf("expected no error, got: %v", gotErr)
			}

			if gotBool != tt.wantBool {
				t.Errorf("got bool %v, want %v", gotBool, tt.wantBool)
			}

			if tt.wantInfo && gotInfo != nil {
				if gotInfo.Name() != filepath.Base(tt.path) {
					t.Errorf("got name %q, want %q", gotInfo.Name(), filepath.Base(tt.path))
				}
			}
		})
	}
}

func TestListEntries(t *testing.T) {
	tmpDir := t.TempDir()

	os.WriteFile(filepath.Join(tmpDir, "file_1.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file_2.txt"), []byte("test"), 0644)
	subDir := filepath.Join(tmpDir, "subpath")
	os.MkdirAll(subDir, 0755)
	os.WriteFile(filepath.Join(subDir, "sub_file_1.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(subDir, "sub_file_2.txt"), []byte("test"), 0644)

	// Expected return
	expectedSuccess := fmt.Sprintf(
		"- file_1.txt (file)\n" +
			"- file_2.txt (file)\n" +
			"- subpath (directory)\n" +
			"  - sub_file_1.txt (file)\n" +
			"  - sub_file_2.txt (file)\n",
	)

	tests := []struct {
		name   string
		path   string
		expect string
		err    error
	}{
		{
			name:   "esisting entries and subentries",
			path:   tmpDir,
			expect: expectedSuccess,
			err:    errors.New(""),
		},
		{
			name:   "passing a file path",
			path:   "/not/exists/dir",
			expect: "path not found at /not/exists/dir",
			err:    errors.New(""),
		},
		{
			name:   "directorie do not exists",
			path:   filepath.Join(tmpDir, "file_1.txt"),
			expect: "path is not a directory",
			err:    errors.New(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entries, err := listEntries(tt.path, 3, "")

			if err != nil {
				if tt.err != errors.New("") && err != tt.err {
					t.Errorf("got unexpected error %v", err)
				}
			}

			if entries != tt.expect {
				t.Errorf("Expected:\n %v\nGot:\n%q", tt.expect, entries)
			}
		})
	}
}
