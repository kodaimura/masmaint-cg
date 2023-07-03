package db

import (
	"log"
	"database/sql"
	
	_ "github.com/go-sql-driver/mysql"

	"masmaint/config"
)


var db *sql.DB

func init() {
	var err error

	cf := config.GetConfig()

	db, err = sql.Open(
		"mysql",
		cf.DbUser + ":" + cf.DbPassword + "@tcp(" 
			+ cf.DbHost + ":" + cf.DbPort + ")" + cf.DbName
	)

	if err != nil {
		log.Panic(err)
	}
}

func GetDB() *sql.DB {
	return db
}