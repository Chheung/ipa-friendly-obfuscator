package utils

import (
	"bytes"
	"fmt"
	"log"
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

func ProcessFile(filePath string, replacements map[string]string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", filePath, err)
	}

	modifiedData := data
	for searchSeq, replaceSeq := range replacements {
		modifiedData = bytes.Replace(modifiedData, []byte(searchSeq), []byte(replaceSeq), -1)
	}

	err = os.WriteFile(filePath, modifiedData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %v", filePath, err)
	}

	fmt.Printf("File %s modified successfully\n", filePath)
	return nil
}

func ProcessDirectory(dirPath string, replacements map[string]string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			err := ProcessFile(path, replacements)
			if err != nil {
				log.Printf("Failed to process file %s: %v", path, err)
			}
		}

		return nil
	})
}

func ValidateReplacements(replacements map[string]string) {
	for key, value := range replacements {
		if len(key) != len(value) {
			log.Printf("replacement length mismatch: key '%s' (length %d) and value '%s' (length %d) do not match", key, len(key), value, len(value))
			panic("WTF")
		}
	}
}