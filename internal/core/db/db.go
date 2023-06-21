package db

import (
	"log"
	"database/sql"
	
	_ "github.com/mattn/go-sqlite3"

	"masmaint/config"
)


var db *sql.DB

func init() {
	var err error

	cf := config.GetConfig()

	dbname := "./" + cf.DbName + ".db"
	db, err = sql.Open("sqlite3", dbname)

	if err != nil {
		log.Panic(err)
	}
}

func GetDB() *sql.DB {
	return db
}