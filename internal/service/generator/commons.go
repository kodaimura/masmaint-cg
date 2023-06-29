package generator

import (
	"io"
	"os"
	"path/filepath"

	"masmaint-cg/internal/core/logger"
)


func WriteFile(path, content string) error {
	f, err := os.Create(path)
	defer f.Close()

	if err != nil {
		logger.LogError(err.Error())
		return err
	}
	if _, err = f.Write([]byte(content)); err != nil {
		logger.LogError(err.Error())
		return err
	}
	return nil
}


func CopyFile(source string, destination string) error {
	src, err := os.Open(source)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return nil
}


func CopyDir(source string, destination string) error {
	err := os.MkdirAll(destination, 0755)
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(source)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(source, entry.Name())
		destinationPath := filepath.Join(destination, entry.Name())

		if entry.IsDir() {
			err := CopyDir(sourcePath, destinationPath)
			if err != nil {
				return err
			}
		} else {
			err := CopyFile(sourcePath, destinationPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}