package service

import (
	"encoding/csv"
	"fmt"
	"os"

	"masmaint/internal/core/logger"
	"masmaint/internal/shared/dto"
)


type CsvParseService struct {}


func NewCsvParseService() *CsvParseService {
	return &CsvParseService{}
}


func (serv *CsvParseService) Parse(path string) ([]dto.Table, error) {
	records, err := serv.readFile(path)
	fmt.Println(records)
	return []dto.Table{}, err
}


func (serv *CsvParseService) readFile(path string) ([][]string, error) {
	file, err := os.Open(path)

	if err != nil {
		logger.LogError("failed to open file:" + path + " " + err.Error())
		return [][]string{}, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	if err != nil {
		logger.LogError("failed to read file:" + path + " " + err.Error())
	}
	return records, err
}