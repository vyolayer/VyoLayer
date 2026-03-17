package mail

import (
	"fmt"
	"net/smtp"
	"strings"
)

// SMTPConfig holds configuration for the SMTP Mailer
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type smtpMailer struct {
	config SMTPConfig
}

// NewSMTPMailer creates a new Mailer that sends emails via SMTP
func NewSMTPMailer(cfg SMTPConfig) Mailer {
	return &smtpMailer{
		config: cfg,
	}
}

func (m *smtpMailer) Send(msg *Message) error {
	auth := smtp.PlainAuth("", m.config.Username, m.config.Password, m.config.Host)

	contentType := "text/plain; charset=\"utf-8\""
	if msg.IsHTML {
		contentType = "text/html; charset=\"utf-8\""
	}

	headers := make(map[string]string)
	headers["From"] = m.config.From
	headers["To"] = strings.Join(msg.To, ",")
	headers["Subject"] = msg.Subject
	headers["Content-Type"] = contentType

	var message strings.Builder
	for k, v := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	message.WriteString("\r\n")
	message.WriteString(msg.Body)

	addr := fmt.Sprintf("%s:%d", m.config.Host, m.config.Port)
	
	err := smtp.SendMail(
		addr,
		auth,
		m.config.From,
		msg.To,
		[]byte(message.String()),
	)
	
	if err != nil {
		return fmt.Errorf("failed to send SMTP email: %w", err)
	}

	return nil
}
