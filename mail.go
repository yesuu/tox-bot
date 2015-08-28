package main

import (
	"fmt"
	"net/smtp"
	"strings"
)

var (
	SmtpUsername string
	SmtpPassword string
	SmtpHost     string
	SmtpAddr     string
)

func Send(subject string, message string, from string, to []string) error {
	auth := smtp.PlainAuth("", SmtpUsername, SmtpPassword, SmtpHost)
	msg := fmt.Sprintf("To: %s\r\nFrom: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		strings.Join(to, ";"), from, subject, message)
	return smtp.SendMail(SmtpAddr, auth, from, to, []byte(msg))
}

func SendHtml(subject string, message string, from string, to []string) error {
	auth := smtp.PlainAuth("", SmtpUsername, SmtpPassword, SmtpHost)
	msg := fmt.Sprintf("To: %s\r\nFrom: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		strings.Join(to, ";"), from, subject, message)
	return smtp.SendMail(SmtpAddr, auth, from, to, []byte(msg))
}
