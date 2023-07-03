package db

import (
	"log"
	"database/sql"
	
	_ "github.com/lib/pq"

	"masmaint/config"
)


var db *sql.DB

func init() {
	var err error

	cf := config.GetConfig()

	db, err = sql.Open(
		"postgres",
		"host=" + cf.DbHost + " port=" + cf.DbPort +
		" user=" + cf.DbUser + " password=" + cf.DbPassword +
		" dbname=" + cf.DbName + " sslmode=disable",
	)

	if err != nil {
		log.Panic(err)
	}
}

func GetDB() *sql.DB {
	return db
}