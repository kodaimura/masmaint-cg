package module

import (
	"reflect"
    "regexp"
    "strings"
	"encoding/json"
	"github.com/lib/pq"
    "github.com/go-sql-driver/mysql"
    "github.com/mattn/go-sqlite3"

	"masmaint/internal/core/errs"
)


func GetConflictColumn(err error) (string, bool) {
    if err == nil {
        return "", false
    }

	if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
        re := regexp.MustCompile(`Key \((\w+)\)=`)
		match := re.FindStringSubmatch(pgErr.Detail)
		if len(match) > 1 {
			return match[1], true
		}
    } else if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
        re := regexp.MustCompile(`for key '([\w.]+)'`)
		match := re.FindStringSubmatch(mysqlErr.Message)
		if len(match) > 1 {
			return match[1], true
		}
    } else if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.Code == sqlite3.ErrConstraint {
        re := regexp.MustCompile(`UNIQUE constraint failed: ([\w.]+)`)
		match := re.FindStringSubmatch(sqliteErr.Error())
		if len(match) > 1 {
			return match[1], true
		}
    }

	return "", false
}

func NewBindError(err error, dataStruct interface{}) error {
    if err == nil {
        return nil
    }

	if jsonErr, ok := err.(*json.UnmarshalTypeError); ok {
		return errs.NewBadRequestError(jsonErr.Field)
	}
	if _, ok := err.(*json.SyntaxError); ok {
		return errs.NewBadRequestError("")
	}
	if strings.Contains(err.Error(), "Key:") {
		fieldName := extractFieldName(err.Error())
		return errs.NewBadRequestError(getFieldJsonTag(dataStruct, fieldName))
	}

    return errs.NewBadRequestError("")
}


func getFieldJsonTag(dataStruct interface{}, fieldName string) string {
    val := reflect.TypeOf(dataStruct).Elem()

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)
        if field.Name == fieldName {
			jsonTag := field.Tag.Get("json")
			if jsonTag != "" && jsonTag != "-" {
				return jsonTag
			}
			return fieldName
		}
    }
    return fieldName
}

func extractFieldName(errorMsg string) string {
	re := regexp.MustCompile(`Key:\s*'([^']+)'`)
	matches := re.FindStringSubmatch(errorMsg)
	if len(matches) > 1 {
		return strings.Split(matches[1], ".")[1]
	}
	return ""
}