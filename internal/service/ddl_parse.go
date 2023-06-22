package service

import (
	"io"
	"os"
	"fmt"
	"masmaint/internal/core/logger"
	"masmaint/internal/shared/dto"
)


type DdlParseService struct {}


func NewDdlParseService() *DdlParseService {
	return &DdlParseService{}
}

func (serv *DdlParseService) readFile(path string) (string, error) {
	ret := ""
	file, err := os.Open(path)

	if err != nil {
		logger.LogError("ファイルオープン失敗:" + path)
		return "", err
	}
	defer file.Close()

	data := make([]byte, 1024)

	for {
		n, err := file.Read(data)
		if err != nil {
			if err == io.EOF {
				break
			}
			logger.LogError("ファイルの読み込み失敗:" + path)
			return "", err
		}
		ret += string(data[:n])
	}
	return ret, nil
}

func (serv *DdlParseService) Parse(path, dbtype string) []dto.Table {
	fmt.Println(serv.readFile(path))
	return []dto.Table{}, err
}
