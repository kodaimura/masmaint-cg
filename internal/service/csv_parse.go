package service

import (
	"encoding/csv"
	"fmt"
	"os"
	//"io"

	"masmaint/internal/core/logger"
	"masmaint/internal/shared/dto"
)


type CsvParseService struct {}


func NewCsvParseService() *CsvParseService {
	return &CsvParseService{}
}


func (serv *CsvParseService) Parse(path string) ([]dto.Table, []string) {
	records, errs := serv.readFile(path)
	if len(errs) != 0 {
		return nil, errs
	}

	errs = serv.validate(records)
	if len(errs) != 0 {
		return nil, errs
	}

	tables := serv.convertTables(records)
	fmt.Println(tables)

	return tables, nil
}


func (serv *CsvParseService) readFile(path string) ([][]string, []string) {
	file, err := os.Open(path)

	if err != nil {
		logger.LogError("failed to open file:" + path + " " + err.Error())
		return nil, []string {"ファイルの読み込みに失敗しました。"}
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()

	if err != nil {
		logger.LogError("failed to read file:" + path + " " + err.Error())
		return nil, []string {"ファイルの読み込みに失敗しました。"}
	}

	return records, nil
}


func (serv *CsvParseService) validate(records [][]string) []string {
	var errs []string
	tcount := 0

	for i, row := range records {

		if row[0] == "t" {
			tcount += 1
			if len(row) < 3 {
				errs = append(errs, fmt.Sprintf("%d行目: 要素数不足", i))
			}
			if row[1] == "" {
				errs = append(errs, fmt.Sprintf("%d行目: テーブル名無し", i))
			}
		}
		if row[0] == "c" {
			if tcount == 0 {
				errs = append(errs, fmt.Sprintf("カラム[%s]のテーブル定義無し", row[1]))
			}
			if len(row) < 6 {
				errs = append(errs, fmt.Sprintf("%d行目: 要素数不足", i))
			}
			if row[1] == "" {
				errs = append(errs, fmt.Sprintf("%d行目: カラム名無し", i))
			}
			if row[3] != "s" && row[3] != "n" {
				errs = append(errs, fmt.Sprintf("%d行目: カラムデータ型に不正な値[%s]", i, row[3]))
			}
			if row[4] != "0" && row[4] != "1" {
				errs = append(errs, fmt.Sprintf("%d行目: NotNull制約フラグに不正な値[%s]", i, row[4]))
			}
			if row[5] != "0" && row[5] != "1" {
				errs = append(errs, fmt.Sprintf("%d行目: 更新不可フラグに不正な値[%s]", i, row[5]))
			}
		}
	}
	return errs
}


func (serv *CsvParseService) convertTables(records [][]string) []dto.Table {
	var tables []dto.Table
	var t dto.Table
	var columns []dto.Column
	var c dto.Column
	tcount := 0

	for _, row := range records {
		if row[0] == "t" {
			if tcount != 0 {
				t.Columns = columns
				tables = append(tables, t)
			}
			tcount += 1
			t = dto.Table{}
			t.TableName = row[1]
			t.TableNameJp = row[2]
			columns = []dto.Column{}
			
		} else if row[0] == "c" {
			c.ColumnName = row[1]
			c.ColumnNameJp = row[2]
			c.ColumnType = row[3]
			c.IsNotNull = (row[4] == "1")
			c.IsReadOnly = (row[5] == "1")
			columns = append(columns, c)
		}
	}
	t.Columns = columns
	tables = append(tables, t)

	return tables
}
