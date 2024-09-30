package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(to string, subject string, body string) error {
	from := os.Getenv("MAIL")
	password := os.Getenv("PASSWORD")
	host := os.Getenv("SMTPHOST")
	port := os.Getenv("SMTPPORT")
	headers := "From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/plain; charset=\"UTF-8\";\r\n\r\n"

	msg := headers + body

	auth := smtp.PlainAuth("", from, password, host)

	err := smtp.SendMail(host+":"+port, auth, from, []string{to}, []byte(msg))
	if err != nil {
		return err
	}

	fmt.Println("Email sent successfully!")
	return nil
}
