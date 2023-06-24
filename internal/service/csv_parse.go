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
	tables, err := serv.convertTables(records)

	fmt.Println(tables)

	return tables, err
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


func (serv *CsvParseService) convertTables(records [][]string) ([]dto.Table, error) {
	var tables []dto.Table
	var t dto.Table
	var columns []dto.Column
	var c dto.Column

	for i, row := range records {
		if row[2] == "" {
			if i != 0 {
				t.Columns = columns
				tables = append(tables, t)
			}
			t = dto.Table{}
			t.TableName = row[0]
			t.TableNameJp = row[1]
			columns = []dto.Column{}
			
		} else {
			c.ColumnName = row[0]
			c.ColumnNameJp = row[1]
			c.ColumnType = row[2]
			c.IsNotNull = (row[3] == "1")
			c.IsReadOnly = (row[4] == "1")
			columns = append(columns, c)
		}
	}
	t.Columns = columns
	tables = append(tables, t)

	return tables, nil
}
