package utils

import (
	"crypto/rand"
	"log"
	"math/big"
	"net/smtp"
	"os"
	"strconv"
)

func SendVerificationEmail(email, code string) {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASS")

	if host == "" || port == "" || user == "" || password == "" {
		log.Printf("SMTP config missing. Mock sending email to %s: code %s\n", email, code)
		return
	}

	auth := smtp.PlainAuth("", user, password, host)
	to := []string{email}
	msg := []byte("To: " + email + "\r\n" +
		"Subject: Verify Your Account\r\n" +
		"\r\n" +
		"Your verification code is: " + code + "\r\n")

	addr := host + ":" + port
	err := smtp.SendMail(addr, auth, user, to, msg)
	if err != nil {
		log.Printf("Error sending email: %v\n", err)
	} else {
		log.Printf("Email sent to %s\n", email)
	}
}

func GenerateVerificationCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}

	code := strconv.FormatInt(n.Int64(), 10)
	for len(code) < 6 {
		code = "0" + code
	}

	return code, nil
}
