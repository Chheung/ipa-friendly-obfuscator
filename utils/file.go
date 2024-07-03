package utils

import (
	"os"
	"path/filepath"
)

func CopyDir(src string, dest string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(dest, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		} else {
			return CopyFile(path, destPath)
		}
	})
}

func CopyFile(src string, dest string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dest, input, 0644)
}

func RemoveDirIfExists(dir string) error {
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return os.RemoveAll(dir)
	}
	return nil
}
