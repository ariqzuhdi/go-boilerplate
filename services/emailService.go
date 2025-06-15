package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
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

func SendVerificationEmail(toEmail, token, recoveryKey string) error {
	link := fmt.Sprintf(os.Getenv("FE_DOMAIN")+"/verify?token=%s", token)
	subject := "Email Verification"

	body := fmt.Sprintf(
		"Please verify your email by clicking the following link:\n\n%s\n\n"+
			"Also, please keep this recovery key safe:\n\n%s\n\n"+
			"Don't share this key with anyone!",
		link, recoveryKey,
	)

	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	msg := []byte("To: " + toEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=\"utf-8\"\r\n" +
		"\r\n" +
		body + "\r\n")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	addr := smtpHost + ":" + smtpPort
	err := smtp.SendMail(addr, auth, from, []string{toEmail}, msg)
	if err != nil {
		log.Printf("SMTP error: %v\n", err)
	} else {
		log.Printf("Verification email sent to %s\n", toEmail)
	}
	return err
}
