package mail

// Message represents an email message to be sent
type Message struct {
	To      []string
	Subject string
	Body    string
	IsHTML  bool
}

// Mailer defines the interface for an email sending service
type Mailer interface {
	Send(msg *Message) error
}
