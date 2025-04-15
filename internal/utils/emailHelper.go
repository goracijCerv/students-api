package utils

import (
	"fmt"
	"net/smtp"

	"github.com/goracijCerv/students-api/internal/config"
)

type EmailHelperSmtp struct {
	smtpAuth smtp.Auth
	hostPort string
	from     string
}

func New(cfg *config.Config) *EmailHelperSmtp {
	auth := smtp.PlainAuth("", cfg.From, cfg.Password, cfg.SmtpHost)
	hostPort := cfg.SmtpHost + ":" + cfg.SmtpPort
	from := cfg.From

	return &EmailHelperSmtp{
		smtpAuth: auth,
		hostPort: hostPort,
		from:     from,
	}
}

func (e *EmailHelperSmtp) SimpleEmailSend(subject, body, to string) error {
	msg := []byte("Subject: " + subject + "\r\n" + "To: " + to + "\r\n" + "Content-Type: text/html; charset=UTF-8\r\n\r\n" + body)
	err := smtp.SendMail(e.hostPort, e.smtpAuth, e.from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("fail to send email: %v", err)
	}
	return nil
}
