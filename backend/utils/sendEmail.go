package utils

import (
	"net/smtp"
	"os"
)

func SendSimpleEmail(to string, subject string, content string) bool {

	auth := smtp.PlainAuth(
		"",
		os.Getenv("SMTP_EMAIL"),
		os.Getenv("EMAIL_PASSWORD"),
		os.Getenv("SMTP"),
	)

	message := "Subject: " + subject + "\n" + content

	err := smtp.SendMail(
		os.Getenv("SMTP")+":"+os.Getenv("SMTP_PORT"),
		auth,
		os.Getenv("SMTP_EMAIL"),
		[]string{to},
		[]byte(message),
	)
	if err != nil {
		return false
	}
	return true
}
