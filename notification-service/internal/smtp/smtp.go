package smtp

import (
	"gopkg.in/gomail.v2"
	"notification-service/config"
)

type EmailSender struct {
	dialer *gomail.Dialer
}

func NewEmailSender(cfg *config.Config) *EmailSender {
	return &EmailSender{
		dialer: gomail.NewDialer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUsername, cfg.SMTPPassword),
	}
}

func (e *EmailSender) SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.dialer.Username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)
	return e.dialer.DialAndSend(m)
}
