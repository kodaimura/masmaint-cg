package generator

import (
	"os"

	"masmaint-cg/internal/core/logger"
)


func WriteFile(path, content string) {
	f, err := os.Create(path)
	defer f.Close()

	if err != nil {
		logger.LogError(err.Error())
	}
	if _, err = f.Write([]byte(content)); err != nil {
		logger.LogError(err.Error())
	}
}