package db

import (
	"fmt"
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
		fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cf.DbHost, cf.DbPort, cf.DbUser, cf.DbPassword, cf.DbName,
		),
	)

	if err != nil {
		log.Panic(err)
	}
}

func GetDB() *sql.DB {
	return db
}