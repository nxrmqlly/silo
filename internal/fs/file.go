package fs

import (
	"os"
	"path/filepath"
)

func ReadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(content), nil

}

func WriteFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

func CreateFile(path string) error {

	//magic flags and permission numbers, I have no idea what the do
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	// this creates the file, by "opening" it, like `touch`
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	return f.Close()
}

func DeletePath(path string) error {
	return os.RemoveAll(path)
}

func RenamePath(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}
