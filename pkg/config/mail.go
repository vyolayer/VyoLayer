package config

import "strconv"

type MailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	UseMock  bool
}

var DefaultMailConfig = MailConfig{
	Host:     "localhost",
	Port:     1025,
	Username: "",
	Password: "",
	From:     "noreply@vyolayer.local",
	UseMock:  true,
}

func NewMailConfig(c MailConfig) *MailConfig {
	return &MailConfig{
		Host:     GetEnv("MAIL_HOST", c.Host),
		Port:     GetEnvInt("MAIL_PORT", strconv.Itoa(c.Port)),
		Username: GetEnv("MAIL_USERNAME", c.Username),
		Password: GetEnv("MAIL_PASSWORD", c.Password),
		From:     GetEnv("MAIL_FROM", c.From),
		UseMock:  GetEnvBool("MAIL_USE_MOCK", strconv.FormatBool(c.UseMock)),
	}
}
