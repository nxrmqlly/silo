package fs_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nxrmqlly/silo/internal/fs"
)

func createTempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "silo-test-*")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	return dir
}

func TestWriteFile(t *testing.T) {
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name     string
		path     string
		content  string
		wantErr  bool
		wantPerm os.FileMode
	}{
		{
			name:     "write simple file",
			path:     filepath.Join(tempDir, "test.txt"),
			content:  "hello world",
			wantErr:  false,
			wantPerm: 0644,
		},
		{
			name:     "write empty file",
			path:     filepath.Join(tempDir, "empty.txt"),
			content:  "",
			wantErr:  false,
			wantPerm: 0644,
		},
		{
			name:     "overwrite existing file",
			path:     filepath.Join(tempDir, "overwrite.txt"),
			content:  "new content",
			wantErr:  false,
			wantPerm: 0644,
		},
		{
			name:    "write to non-existent directory",
			path:    filepath.Join(tempDir, "nonexistent", "file.txt"),
			content: "test",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fs.WriteFile(tt.path, tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify file was written
				data, err := os.ReadFile(tt.path)
				if err != nil {
					t.Errorf("failed to read written file: %v", err)
					return
				}

				if string(data) != tt.content {
					t.Errorf("file content = %q, want %q", string(data), tt.content)
				}

				// Verify permissions
				info, err := os.Stat(tt.path)
				if err != nil {
					t.Errorf("failed to stat file: %v", err)
					return
				}

				if info.Mode().Perm() != tt.wantPerm {
					t.Errorf("file permissions = %o, want %o", info.Mode().Perm(), tt.wantPerm)
				}
			}
		})
	}
}

func TestReadFile(t *testing.T) {
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	// Setup test files
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "hello from test file"
	os.WriteFile(testFile, []byte(testContent), 0644)

	tests := []struct {
		name    string
		path    string
		want    string
		wantErr bool
	}{
		{
			name:    "read existing file",
			path:    testFile,
			want:    testContent,
			wantErr: false,
		},
		{
			name:    "read non-existent file",
			path:    filepath.Join(tempDir, "nonexistent.txt"),
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fs.ReadFile(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("ReadFile() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCreateFile(t *testing.T) {
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "create file in existing directory",
			path:    filepath.Join(tempDir, "newfile.txt"),
			wantErr: false,
		},
		{
			name:    "create file with parent directories",
			path:    filepath.Join(tempDir, "subdir1", "subdir2", "newfile.txt"),
			wantErr: false,
		},
		{
			name:    "create file that already exists",
			path:    filepath.Join(tempDir, "existing.txt"),
			wantErr: true,
		},
	}

	// pre create one file for "already exists" test
	existingFile := filepath.Join(tempDir, "existing.txt")
	os.WriteFile(existingFile, []byte("exists"), 0644)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fs.CreateFile(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify file exists
				if _, err := os.Stat(tt.path); err != nil {
					t.Errorf("file was not created: %v", err)
				}
			}
		})
	}
}

func TestDeletePath(t *testing.T) {
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name    string
		setup   func(string) string // returns path to delete
		wantErr bool
	}{
		{
			name: "delete single file",
			setup: func(dir string) string {
				path := filepath.Join(dir, "file.txt")
				os.WriteFile(path, []byte("test"), 0644)
				return path
			},
			wantErr: false,
		},
		{
			name: "delete empty directory",
			setup: func(dir string) string {
				path := filepath.Join(dir, "emptydir")
				os.Mkdir(path, 0755)
				return path
			},
			wantErr: false,
		},
		{
			name: "delete directory with files",
			setup: func(dir string) string {
				path := filepath.Join(dir, "dirwithfiles")
				os.Mkdir(path, 0755)
				os.WriteFile(filepath.Join(path, "file1.txt"), []byte("test1"), 0644)
				os.WriteFile(filepath.Join(path, "file2.txt"), []byte("test2"), 0644)
				return path
			},
			wantErr: false,
		},
		{
			name: "delete non-existent path",
			setup: func(dir string) string {
				return filepath.Join(dir, "nonexistent")
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create fresh temp directory for each test
			testDir := createTempDir(t)
			defer os.RemoveAll(testDir)

			pathToDelete := tt.setup(testDir)
			err := fs.DeletePath(pathToDelete)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeletePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify path was deleted
				if _, err := os.Stat(pathToDelete); !os.IsNotExist(err) {
					t.Errorf("path was not deleted")
				}
			}
		})
	}
}
