package db

import (
	"log"
	"fmt"
	"reflect"
	"strings"	
	"database/sql"
	
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"masmaint/config"
	"masmaint/internal/core/utils"
)


var db *sql.DB
var driver string

func init() {
	var err error

	cf := config.GetConfig()
	driver = cf.DBDriver
	var dsn string

	if driver == "sqlite3" {
		dsn = cf.DBName
	} else if driver == "mysql" {
		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s", 
			cf.DBUser, cf.DBPass, cf.DBHost, cf.DBPort, cf.DBName,
		)
	} else if driver == "postgres" {
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cf.DBHost, cf.DBPort, cf.DBUser, cf.DBPass, cf.DBName,
		)
	} else {
		log.Panic("Error: must specify a valid DB_DRIVER: 'postgres', 'mysql', or 'sqlite3'.")
	}

	db, err = sql.Open(driver, dsn)
	if err != nil {
		log.Panic(err)
	}
}

func GetDB() *sql.DB {
	return db
}

func getBindVar (seq int) string {
	if driver == "postgres" {
		return fmt.Sprintf("$%d", seq)
	} else {
		return "?"
	}
}

func BuildWhereClause(filter interface{}) (string, []interface{}) {
	var conditions []string
	var binds []interface{}

	val := reflect.ValueOf(filter)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	seq := 1
	for i := 0; i < val.NumField(); i++ {
		columnName := typ.Field(i).Tag.Get("db")
		fieldValue := val.Field(i).Interface()

		if columnName == "" { 
			continue
		}
		if utils.IsZero(fieldValue) {
			continue
		}

		conditions = append(conditions, fmt.Sprintf("%s = %s", columnName, getBindVar(seq)))
		binds = append(binds, fieldValue)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	return whereClause, binds
}