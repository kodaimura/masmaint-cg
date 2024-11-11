package mail

import (
	"fmt"
	"strings"
    "net/smtp"
	"mime"
	
	"masmaint/config"
)


func Send (from string, to []string, subject, body string) error {
	address := getSmtpAddress()
	auth := getSmtpAuth()
	msg := generateMessageByte(from, to, subject, body)
	return smtp.SendMail(address, auth, from, to, msg)
}

func getSmtpAuth () smtp.Auth {
	cf := config.GetConfig()
	return smtp.PlainAuth("", cf.MailUser, cf.MailPass, cf.MailHost)
}

func getSmtpAddress () string {
	cf := config.GetConfig()
	return fmt.Sprintf("%s:%s", cf.MailHost, cf.MailPort)
}

func generateMessageByte (from string, to []string, subject, body string) []byte {
	header := "From: " + from + "\r\n" +
		"To: " + strings.Join(to, ", ") + "\r\n" +
		"Subject: " + mime.QEncoding.Encode("UTF-8", subject) + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n"

	return []byte(header + "\r\n" + body)
}