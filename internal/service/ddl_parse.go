package service

import (
	"io"
	"os"
	"strings"
	"errors"

	"masmaint-cg/internal/core/logger"
	"masmaint-cg/internal/shared/dto"
)


type ddlParseService struct {}


func NewDdlParseService() *ddlParseService {
	return &ddlParseService{}
}


func (serv *ddlParseService) Parse(path, dbtype string) ([]dto.Table, error) {
	ddl, _ := serv.readFile(path)
	tokens := serv.lexicalAnalysis(ddl)
	err := serv.validate(tokens, dbtype)
	return []dto.Table{}, err
}


func (serv *ddlParseService) readFile(path string) (string, error) {
	ret := ""
	file, err := os.Open(path)

	if err != nil {
		logger.LogError("failed to open file:" + path)
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
			logger.LogError("failed to read file:" + path)
			return "", err
		}
		ret += string(data[:n])
	}
	return ret, nil
}


//字句解析
func (serv *ddlParseService) lexicalAnalysis(ddl string) []string {
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

func (serv *ddlParseService) validate(tokens []string, dbtype string) error {
	if (dbtype == "postgresql") {
		return serv.validatePostgreSQL(tokens, dbtype)
	} else {
		return errors.New("not supported")
	}
}

func (serv *ddlParseService) validatePostgreSQL(tokens []string, dbtype string) error {
	return nil
}