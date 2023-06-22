package service

import (
	"io"
	"os"
	"fmt"
	"strings"
	"masmaint/internal/core/logger"
	"masmaint/internal/shared/dto"
)


type DdlParseService struct {}


func NewDdlParseService() *DdlParseService {
	return &DdlParseService{}
}


func (serv *DdlParseService) Parse(path, dbtype string) []dto.Table {
	ddl, _ := serv.readFile(path)
	tokens := serv.LexicalAnalysis(ddl)
	fmt.Println(tokens)
	return []dto.Table{}
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


//字句解析
func (serv *DdlParseService) LexicalAnalysis(ddl string) []string {
	ddl = strings.ReplaceAll(ddl, "\n", "")
	ddl = strings.ReplaceAll(ddl, "(", " ( ")
	ddl = strings.ReplaceAll(ddl, ")", " ) ")
	ddl = strings.ReplaceAll(ddl, ",", " , ")
	ddl = strings.ReplaceAll(ddl, ";", " ; ")
	ddl = strings.ReplaceAll(ddl, "'", " ' ")
	ddl = strings.ReplaceAll(ddl, "\"", " \" ")
	ddl = strings.ReplaceAll(ddl, " ` ", " ` ")

	return strings.Split(ddl, " ")
} 