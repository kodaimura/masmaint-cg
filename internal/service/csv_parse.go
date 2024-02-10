package service

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"masmaint-cg/internal/core/logger"
	"masmaint-cg/internal/shared/dto"
)


type csvParseService struct {}


func NewCsvParseService() *csvParseService {
	return &csvParseService{}
}


func (serv *csvParseService) Parse(path string) ([]dto.Table, []string) {
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


func (serv *csvParseService) readFile(path string) ([][]string, []string) {
	file, err := os.Open(path)

	if err != nil {
		logger.Error("failed to open file:" + path + " " + err.Error())
		return nil, []string {"ファイルの読み込みに失敗しました。"}
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()

	if err != nil {
		logger.Error("failed to read file:" + path + " " + err.Error())
		return nil, []string {"ファイルの読み込みに失敗しました。"}
	}

	return records, nil
}

func (serv *csvParseService) isValidTableLabel(s string) bool {
	return s == "t" || s == "T"
}

func (serv *csvParseService) isValidColumnLabel(s string) bool {
	return s == "c" || s == "C"
}

func (serv *csvParseService) isValidColumnName(s string) bool {
	//一旦空文字でないかのチェックのみ
	//利用可能文字チェックなど実装想定
	return s != ""
}

func (serv *csvParseService) isValidTableName(s string) bool {
	//一旦空文字でないかのチェックのみ
	//利用可能文字チェックなど実装想定
	return s != ""
}

func (serv *csvParseService) isValidColumnType(s string) bool {
	sl := strings.ToLower(s)
	return sl == "s" || sl == "i" || sl == "f" || sl == "t"
}

func (serv *csvParseService) isValidFlg(s string) bool {
	return s == "0" || s == "1"
}

func (serv *csvParseService) isStartWithTableLabel(records [][]string) bool {
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

func (serv *csvParseService) validate(records [][]string) []string {
	var errs []string
	tname := ""
	hasPk := true

	if !serv.isStartWithTableLabel(records) {
		errs = append(errs, "最初のテーブルが見つかりません。")
	}

	for index, row := range records {
		i := index + 1 //エラーメッセージ用

		if serv.isValidTableLabel(row[0]) {
			if len(row) < 3 {
				errs = append(errs, fmt.Sprintf("%d行: 要素が不足しています。", i))
				continue
			}
			if !serv.isValidTableName(row[1]) {
				errs = append(errs, fmt.Sprintf("%d行2列: テーブル名が不正です。[%s]", i, row[1]))
			}
			if !hasPk {
				errs = append(errs, fmt.Sprintf("主キーの無いテーブルは取り込めません。[%s]", tname))
			} 
			tname = row[1]
			hasPk = false
		}

		if serv.isValidColumnLabel(row[0]) {
			if len(row) < 7 {
				errs = append(errs, fmt.Sprintf("%d行目: 要素が不足しています。", i))
				continue
			}
			if !serv.isValidColumnName(row[1]) {
				errs = append(errs, fmt.Sprintf("%d行2列: カラム名が不正です。[%s]", i, row[1]))
			}
			if !serv.isValidColumnType(row[2]) {
				errs = append(errs, fmt.Sprintf("%d行3列: データ型は I, F, S, T のいずれかとしてください。[%s]", i, row[2]))
			}
			if !serv.isValidFlg(row[3]) {
				errs = append(errs, fmt.Sprintf("%d行4列: 主キーフラグは 0 または 1 としてください。[%s]", i, row[3]))
			}
			if !serv.isValidFlg(row[4]) {
				errs = append(errs, fmt.Sprintf("%d行5列: NotNull制約フラグは 0 または 1 としてください。[%s]", i, row[4]))
			}
			if !serv.isValidFlg(row[5]) {
				errs = append(errs, fmt.Sprintf("%d行6列: 登録可能フラグは 0 または 1 としてください。[%s]", i, row[5]))
			}
			if !serv.isValidFlg(row[6]) {
				errs = append(errs, fmt.Sprintf("%d行7列: 更新可能フラグは 0 または 1 としてください。[%s]", i, row[6]))
			}
			if row[3] == "1" && row[6] != "0" {
				errs = append(errs, fmt.Sprintf("%d行7列: 主キーのカラムは更新不可 0 としてください。", i))
			} 
			if row[3] == "1" && row[4] != "0" {
				errs = append(errs, fmt.Sprintf("%d行5列: 主キーのカラムはNotNull制約フラグを 0 としてください。", i))
			} 
			if row[3] == "1" {
				hasPk = true
			}
		}
	}

	if !hasPk {
		errs = append(errs, fmt.Sprintf("主キーの無いテーブルは取り込めません。[%s]", tname))
	} 

	return errs
}


func (serv *csvParseService) convertTables(records [][]string) []dto.Table {
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
			c.ColumnType = strings.ToLower(row[2])
			c.IsPrimaryKey = (row[3] == "1")
			c.IsNotNull = (row[4] == "1")
			c.IsInsAble = (row[5] == "1")
			c.IsUpdAble = (row[6] == "1")
			columns = append(columns, c)
		}
	}
	t.Columns = columns
	tables = append(tables, t)

	return tables
}
