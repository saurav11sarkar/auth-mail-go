package email

import (
	"fmt"
	"mime"
	"net/mail"
	"net/smtp"
	"strings"

	"github.com/saurav11sarkar/go/internal/config"
)

type Email struct {
	cfg config.Config
}

func NewEmail(cfg config.Config) *Email {
	return &Email{cfg: cfg}
}

func (e *Email) Send(to, subject, body string) error {
	from, err := mail.ParseAddress(e.cfg.Smtp.Username)
	if err != nil {
		return fmt.Errorf("invalid sender email %w", err)
	}

	toEmail, err := mail.ParseAddress(to)
	if err != nil {
		return fmt.Errorf("invalid recipient email %w", err)
	}

	if strings.ContainsAny(subject, "\r\n") {
		return fmt.Errorf("invalid subject: contains newline characters")
	}

	header := map[string]string{
		"From":         from.String(),
		"To":           toEmail.String(),
		"Subject":      mime.QEncoding.Encode("utf-8", subject),
		"Content-Type": "text/html; charset=UTF-8",
		"MIME-Version": "1.0",
	}

	var msg strings.Builder

	for k, v := range header {
		fmt.Fprintf(&msg, "%s: %s\r\n", k, v)
	}
	msg.WriteString("\r\n")
	msg.WriteString(body)

	auth := smtp.PlainAuth("", e.cfg.Smtp.Username, e.cfg.Smtp.Password, e.cfg.Smtp.Host)

	if err := smtp.SendMail(e.cfg.Smtp.Host+":"+e.cfg.Smtp.Port, auth, from.Address, []string{toEmail.Address}, []byte(msg.String())); err != nil {
		return fmt.Errorf("send email failed %w", err)
	}
	return nil

}
