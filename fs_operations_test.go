package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
		name          string
		path          string
		expectContent string
		expectMessage string
		err           error
	}{
		{
			name:          "esisting entries and subentries",
			path:          tmpDir,
			expectContent: expectedSuccess,
			expectMessage: "",
			err:           errors.New(""),
		},
		{
			name:          "passing a file path",
			path:          "/not/exists/dir",
			expectContent: "",
			expectMessage: "path not found at /not/exists/dir",
			err:           errors.New(""),
		},
		{
			name:          "directorie do not exists",
			path:          filepath.Join(tmpDir, "file_1.txt"),
			expectContent: "",
			expectMessage: "path is not a directory",
			err:           errors.New(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			operationResult := listEntries(tt.path, 3, "")

			if operationResult.Error != nil {
				if tt.err != errors.New("") && operationResult.Error != tt.err {
					t.Errorf("got unexpected error %v", operationResult.Error)
				}
			}

			if operationResult.Content != tt.expectContent {
				t.Errorf("Expected:\n %v\nGot:\n%q", tt.expectContent, operationResult.Content)
			}

			if operationResult.Message != tt.expectMessage {
				t.Errorf("Expected:\n %v\nGot:\n%q", tt.expectMessage, operationResult.Message)
			}
		})
	}
}

func TestReadFile(t *testing.T) {
	tmpDir := t.TempDir()

	os.WriteFile(filepath.Join(tmpDir, "file_1.txt"), []byte("test"), 0644)
	subDir := filepath.Join(tmpDir, "subpath")
	os.MkdirAll(subDir, 0755)

	tests := []struct {
		name          string
		path          string
		expectMessage string
		expectContent string
	}{
		{
			name:          "read file sucessfully",
			path:          filepath.Join(tmpDir, "file_1.txt"),
			expectMessage: "",
			expectContent: "test",
		},
		{
			name:          "read file sucessfully",
			path:          subDir,
			expectMessage: "path is a directory, must be a file",
			expectContent: "",
		},
		{
			name:          "read file sucessfully",
			path:          "/not/exists/file.txt",
			expectMessage: "path not found at /not/exists/file.txt",
			expectContent: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			operationResult := readFile(tt.path)
			if operationResult.Error != nil {
				t.Errorf("unexpected error: %v", operationResult.Error)
			}

			if operationResult.Message != tt.expectMessage {
				t.Errorf("Got %s, expected: %s", operationResult.Message, tt.expectMessage)
			}
			if operationResult.Content != tt.expectContent {
				t.Errorf("Got %s, expected: %s", operationResult.Content, tt.expectContent)
			}
		})
	}
}

func TestWriteToFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup test file and directory
	filePath := filepath.Join(tmpDir, "existing.txt")
	if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	tests := []struct {
		name     string
		content  string
		path     string
		expect   string
		contains bool
		isError  bool
	}{
		{
			name:     "write to new file",
			content:  "hello",
			path:     filepath.Join(tmpDir, "newfile.txt"),
			expect:   "file written successfully",
			contains: false,
			isError:  false,
		},
		{
			name:     "overwrite existing file",
			content:  "updated",
			path:     filePath,
			expect:   "file written successfully",
			contains: false,
			isError:  false,
		},
		{
			name:     "path is directory",
			content:  "should fail",
			path:     subDir,
			expect:   "path is a directory, must be a file",
			contains: false,
			isError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := writeToFile(tt.content, tt.path)

			if tt.isError {
				if result.Error == nil {
					t.Fatal("Expected error but got none")
				}
				if !strings.Contains(result.Error.Error(), tt.expect) {
					t.Errorf("Error does not contain expected message.\nGot: %v\nExpected to contain: %s", result.Error, tt.expect)
				}
			} else {
				if result.Error != nil {
					t.Fatalf("Unexpected error: %v", result.Error)
				}

				if tt.contains {
					if !strings.Contains(result.Content, tt.expect) {
						t.Errorf("Output does not contain expected substring.\nGot: %s\nExpected to contain: %s", result.Content, tt.expect)
					}
				} else {
					if result.Content != tt.expect && result.Message != tt.expect {
						t.Errorf("Unexpected result.\nGot Content: %s\nGot Message: %s\nExpected: %s",
							result.Content, result.Message, tt.expect)
					}
				}
			}
		})
	}
}

func TestGetFileInfo(t *testing.T) {
	tmpDir := t.TempDir()

	filePath := filepath.Join(tmpDir, "file_1.txt")
	if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	subDir := filepath.Join(tmpDir, "subpath")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	tests := []struct {
		name     string
		path     string
		expect   string
		contains bool
	}{
		{
			name:     "read file successfully",
			path:     filePath,
			expect:   "File: " + filePath,
			contains: true,
		},
		{
			name:     "path is directory",
			path:     subDir,
			expect:   "path is a directory, must be a file",
			contains: false,
		},
		{
			name:     "path does not exist",
			path:     "/not/exists/file.txt",
			expect:   "path not found at /not/exists/file.txt",
			contains: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := getFileInfo(tt.path)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.contains {
				if !strings.Contains(content, tt.expect) {
					t.Errorf("Output does not contain expected substring.\nGot:\n%s\nExpected to contain:\n%s", content, tt.expect)
				}
			} else {
				if content != tt.expect {
					t.Errorf("Got:\n%s\nExpected:\n%s", content, tt.expect)
				}
			}
		})
	}
}
func TestRenameFilaAndDir(t *testing.T) {
	tmpDir := t.TempDir()

	filePathOne := filepath.Join(tmpDir, "file_1.txt")
	if err := os.WriteFile(filePathOne, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	filePathTwo := filepath.Join(tmpDir, "file_2.txt")
	if err := os.WriteFile(filePathTwo, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	subDir := filepath.Join(tmpDir, "subpath")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	tests := []struct {
		name        string
		path        string
		newPathName string
		expect      string
	}{
		{
			name:        "rename a file",
			path:        filePathOne,
			newPathName: "updated_file_name.txt",
			expect:      filepath.Join(tmpDir, "updated_file_name.txt"),
		},
		{
			name:        "rename a directory",
			path:        subDir,
			newPathName: "updated_dir_name",
			expect:      filepath.Join(tmpDir, "updated_dir_name"),
		},
		{
			name:        "path does not exist",
			path:        "/not/exists/file.txt",
			newPathName: "updated_file_name.txt",
			expect:      "path not found at /not/exists/file.txt",
		},
		{
			name:        "file already exists",
			path:        filePathTwo,
			newPathName: "updated_file_name.txt",
			expect:      fmt.Sprintf("target path %s already exists", filepath.Join(tmpDir, "updated_file_name.txt")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := renamePath(tt.path, tt.newPathName)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if content != tt.expect {
				t.Errorf("Got:\n%s\nExpected:\n%s", content, tt.expect)
			}
		})
	}
}

func TestCopyFileOrDir(t *testing.T) {
	tmpDir := t.TempDir()

	filePath := filepath.Join(tmpDir, "file.txt")
	content := []byte("Hello, World!")
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	destinationFilePath := filepath.Join(tmpDir, "file_copy.txt")

	folderPath := filepath.Join(tmpDir, "folder")
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	copyDestinationDirPath := filepath.Join(tmpDir, "folder_copy")

	tests := []struct {
		name        string
		source      string
		destination string
		expect      string
	}{
		{
			name:        "copy file",
			source:      filePath,
			destination: destinationFilePath,
			expect:      "File copied to destination",
		},
		{
			name:        "copy directory",
			source:      folderPath,
			destination: copyDestinationDirPath,
			expect:      "",
		},
		{
			name:        "copy non-existent file",
			source:      filepath.Join(tmpDir, "nonexistent.txt"),
			destination: filepath.Join(tmpDir, "nonexistent_copy.txt"),
			expect:      "path not found at " + filepath.Join(tmpDir, "nonexistent.txt"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := copyFileOrDir(tt.source, tt.destination)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if result != tt.expect {
				t.Errorf("Got: %v, Expected: %v", result, tt.expect)
			}
		})
	}
}
