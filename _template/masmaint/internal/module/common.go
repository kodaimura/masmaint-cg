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
        if pgErr.Detail != "" {
			parts := strings.Split(pgErr.Detail, "violates unique constraint")
			if len(parts) > 0 {
				return strings.TrimSpace(parts[0]), true
			}
		}
    } else if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
        if strings.Contains(mysqlErr.Message, "for key") {
			parts := strings.Split(mysqlErr.Message, "for key")
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1]), true
			}
		}
    } else if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.Code == sqlite3.ErrConstraint {
        if strings.Contains(sqliteErr.Error(), "UNIQUE constraint failed") {
			parts := strings.Split(sqliteErr.Error(), ":")
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1]), true
			}
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