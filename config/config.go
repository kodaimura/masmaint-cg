package config

import (
	"os"
	"log"
	"fmt"

	"github.com/joho/godotenv"
)


type Config struct {
	AppName string
	AppHost string
	AppPort string

	DBDriver string
	DBName string
	DBHost string
	DBPort string
	DBUser string
	DBPass string

	MailHost string
	MailPort string
	MailUser string
	MailPass string

	BasicAuthUser string
	BasicAuthPass string

	JwtSecretKey string
	LogLevel string
}

var cf Config


func init() {
	err := godotenv.Load(fmt.Sprintf("config/env/%s.env", os.Getenv("ENV")))

	if err != nil {
		log.Panic(err)
	}

	cf.AppName = os.Getenv("APP_NAME")
	cf.AppHost = os.Getenv("APP_HOST")
	cf.AppPort = os.Getenv("APP_PORT")

	cf.DBDriver = os.Getenv("DB_DRIVER")
	cf.DBName = os.Getenv("DB_NAME")
	cf.DBHost = os.Getenv("DB_HOST")
	cf.DBPort = os.Getenv("DB_PORT")
	cf.DBUser = os.Getenv("DB_USER")
	cf.DBPass = os.Getenv("DB_PASSWORD")

	cf.MailHost = os.Getenv("MAIL_HOST")
	cf.MailPort = os.Getenv("MAIL_PORT")
	cf.MailUser = os.Getenv("MAIL_USER")
	cf.MailPass = os.Getenv("MAIL_PASSWORD")

	cf.BasicAuthUser = os.Getenv("BASIC_AUTH_USER")
	cf.BasicAuthPass = os.Getenv("BASIC_AUTH_PASSWORD")

	cf.JwtSecretKey = os.Getenv("JWT_SECRET_KEY")
	cf.LogLevel = os.Getenv("LOG_LEVEL")
}


func GetConfig() *Config{
	return &cf
}