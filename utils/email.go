package utils

import (
	"strconv"

	"gopkg.in/gomail.v2"

	"os"
)

func SendMail(	to string, subject string, body string) error {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))


	mailer := gomail.NewMessage()
	mailer.SetHeader("From", os.Getenv("SMTP_FROM"))
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	dialer := gomail.NewDialer(
		"smtp.hostinger.com",
		port,
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASS"),

	)

	return dialer.DialAndSend(mailer)
}

