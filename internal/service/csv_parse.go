package service

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"masmaint-cg/internal/core/logger"
	"masmaint-cg/internal/shared/dto"
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

func (serv *CsvParseService) isValidTableLabel(s string) bool {
	return s == "t" || s == "T"
}

func (serv *CsvParseService) isValidColumnLabel(s string) bool {
	return s == "c" || s == "C"
}

func (serv *CsvParseService) isValidColumnName(s string) bool {
	//一旦空文字でないかのチェックのみ
	//利用可能文字チェックなど実装想定
	return s != ""
}

func (serv *CsvParseService) isValidTableName(s string) bool {
	//一旦空文字でないかのチェックのみ
	//利用可能文字チェックなど実装想定
	return s != ""
}

func (serv *CsvParseService) isValidColumnType(s string) bool {
	sl := strings.ToLower(s)
	return sl == "s" || sl == "i" || sl == "f" || sl == "t"
}

func (serv *CsvParseService) isValidFlg(s string) bool {
	return s == "0" || s == "1"
}

func (serv *CsvParseService) isStartWithTableLabel(records [][]string) bool {
	for _, row := range records {
		if serv.isValidTableLabel(row[0]) {
			return true
		}
		if serv.isValidColumnLabel(row[0]) {
			return false
		}
	}
	return false
}

func (serv *CsvParseService) validate(records [][]string) []string {
	var errs []string

	if !serv.isStartWithTableLabel(records) {
		errs = append(errs, "最初のテーブル定義無し")
	}

	for i, row := range records {

		if serv.isValidTableLabel(row[0]) {
			if len(row) < 3 {
				errs = append(errs, fmt.Sprintf("%d行目: 要素数不足", i))
				continue
			}
			if !serv.isValidTableName(row[1]) {
				errs = append(errs, fmt.Sprintf("%d行目: テーブル名不正[%s]", i, row[1]))
			}
		}
		if serv.isValidColumnLabel(row[0]) {
			if len(row) < 7 {
				errs = append(errs, fmt.Sprintf("%d行目: 要素数不足", i))
				continue
			}
			if !serv.isValidColumnName(row[1]) {
				errs = append(errs, fmt.Sprintf("%d行目: カラム名不正[%s]", i, row[1]))
			}
			if !serv.isValidColumnType(row[3]) {
				errs = append(errs, fmt.Sprintf("%d行目: カラムデータ型不正[%s]", i, row[3]))
			}
			if !serv.isValidFlg(row[4]) {
				errs = append(errs, fmt.Sprintf("%d行目: 主キーフラグ不正[%s]", i, row[4]))
			}
			if !serv.isValidFlg(row[5]) {
				errs = append(errs, fmt.Sprintf("%d行目: NotNull制約フラグ不正[%s]", i, row[5]))
			}
			if !serv.isValidFlg(row[6]) {
				errs = append(errs, fmt.Sprintf("%d行目: 更新不可フラグ不正[%s]", i, row[6]))
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
		if serv.isValidTableLabel(row[0]) {
			if tcount != 0 {
				t.Columns = columns
				tables = append(tables, t)
			}
			tcount += 1
			t = dto.Table{}
			t.TableName = row[1]
			t.TableNameJp = row[2]
			columns = []dto.Column{}
			
		} else if serv.isValidColumnLabel(row[0]) {
			c.ColumnName = row[1]
			c.ColumnNameJp = row[2]
			c.ColumnType = strings.ToLower(row[3])
			c.IsPrimaryKey = (row[4] == "1")
			c.IsNotNull = (row[5] == "1")
			c.IsReadOnly = (row[6] == "1")
			columns = append(columns, c)
		}
	}
	t.Columns = columns
	tables = append(tables, t)

	return tables
}
