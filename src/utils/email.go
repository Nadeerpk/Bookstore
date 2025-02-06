package utils

import (
	"log"
	"net/smtp"
	"os"
)

func Send_mail(receipients []string, subject string, body string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	senderEmail := "nadeer@qburst.com"
	password := os.Getenv("app_password")

	message := []byte(subject + "\n" + body)

	auth := smtp.PlainAuth("", senderEmail, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, receipients, message)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}
	return nil
}
