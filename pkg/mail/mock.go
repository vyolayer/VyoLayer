package mail

import (
	"log"
)

type mockMailer struct{}

// NewMockMailer creates a dummy Mailer that just logs to console
// Useful for development and testing without spamming real addresses
func NewMockMailer() Mailer {
	return &mockMailer{}
}

func (m *mockMailer) Send(msg *Message) error {
	log.Printf(
		"📧 [MOCK MAIL] Sending to %v | Subject: %s\n | Body: %s\n",
		msg.To,
		msg.Subject,
		msg.Body,
	)
	return nil
}
