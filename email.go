package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

const smtpHost = "smtp.gmail.com"
const smtpPort = "587"

func EmailOK(name string) {
	subject := fmt.Sprintf("[ok] %s deployed", name)
	body := fmt.Sprintf("[ok] %s deployed", name)
	sendEmail(subject, body)
}

func EmailERR(name, script, output string, cause error) {
	subject := fmt.Sprintf("[error] %s/%s failed", name, script)
	body := fmt.Sprintf("Cause: %v\nOutput:\n%s", cause, output)
	sendEmail(subject, body)
}

func sendEmail(subject, body string) {
	from := os.Getenv("GDEP_EMAIL_FROM")
	password := os.Getenv("GDEP_EMAIL_PASSWORD")
	to := os.Getenv("GDEP_EMAIL_TO")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s\r\n", from, to, subject, body)

	if err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg)); err != nil {
		log.Println("send email failed:", err)
	}
}
