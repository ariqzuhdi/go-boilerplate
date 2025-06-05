package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/smtp"
	"os"
)

func GenerateToken(n int) (string, error) {
	// Create a new token object
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func SendVerificationEmail(toEmail, token string) error {
	link := fmt.Sprintf("http://localhost:8080/verify?token=%s", token)
	subject := "Email Verification"
	body := fmt.Sprintf("Please verify your email by clicking the following link: %s", link)

	// Here you would send the email using your preferred email service
	// For example:
	// return emailService.SendEmail(toEmail, subject, body)

	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	msg := []byte("Subject: " + subject + "\r\n\r\n" + body)

	auth := smtp.PlainAuth("", from, password, smtpHost)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{toEmail}, msg)
}
