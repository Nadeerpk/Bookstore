package utils

import (
	"log"
	"net/smtp"
	"os"
)

func Send_mail(receipients []string, subject string, body string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	senderEmail := os.Getenv("sender_email")
	password := os.Getenv("app.password")

	message := []byte(subject + "\n" + body)

	auth := smtp.PlainAuth("", senderEmail, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, receipients, message)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}
	return nil
}
