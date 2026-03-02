package fs

import "os"

func CreateDir(path string) error {
	return os.MkdirAll(path, 0755)
}
