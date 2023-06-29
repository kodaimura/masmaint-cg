package generator

import (
	"os"

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