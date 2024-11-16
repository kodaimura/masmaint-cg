package generator

import (
	"io"
	"os"
	"strings"
	"path/filepath"

	"masmaint-cg/internal/core/logger"
)


func WriteFile(path, content string) error {
	f, err := os.Create(path)
	defer f.Close()

	if err != nil {
		logger.Error(err.Error())
		return err
	}
	if _, err = f.WriteString(content); err != nil {
		logger.Error(err.Error())
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


func ReadFile(path string) string {
	file, err := os.Open(path)
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		logger.Fatal(err.Error())
	}
	fileSize := fileInfo.Size()

	content := make([]byte, fileSize)

	_, err = file.Read(content)
	if err != nil {
		logger.Fatal(err.Error())
	}

	return string(content)
}


//xxx -> Xxx / xxx_yyy -> XxxYyy
func SnakeToPascal(snake string) string {
	ls := strings.Split(strings.ToLower(snake), "_")
	for i, s := range ls {
		ls[i] = strings.ToUpper(s[0:1]) + s[1:]
	}
	return strings.Join(ls, "")
}

//xxx -> xxx / xxx_yyy -> xxxYyy
func SnakeToCamel(snake string) string {
	ls := strings.Split(strings.ToLower(snake), "_")
	for i, s := range ls {
		if i != 0 {
			ls[i] = strings.ToUpper(s[0:1]) + s[1:]
		}
	}
	return strings.Join(ls, "")
}

//xxx -> x / xxx_yyy -> xy
func GetSnakeInitial(snake string) string {
	ls := strings.Split(strings.ToLower(snake), "_")
	ret := ""
	for _, s := range ls {
		ret += s[0:1]
	}
	return ret
}

func Contains(slice []string, element string) bool {
    for _, v := range slice {
        if v == element {
            return true
        }
    }
    return false
}