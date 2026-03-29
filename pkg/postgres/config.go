package postgres

import "github.com/vyolayer/vyolayer/pkg/config"

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

var DefaultConfig = Config{
	Host:     "localhost",
	Port:     "4444",
	User:     "vyolayer_user",
	Password: "vyolayer_password",
	DBName:   "vyolayer_db",
	SSLMode:  "disable",
}

func NewConfig(c Config) *Config {
	return &Config{
		Host:     config.GetEnv("DB_HOST", c.Host),
		Port:     config.GetEnv("DB_PORT", c.Port),
		User:     config.GetEnv("DB_USER", c.User),
		Password: config.GetEnv("DB_PASSWORD", c.Password),
		DBName:   config.GetEnv("DB_NAME", c.DBName),
		SSLMode:  config.GetEnv("DB_SSLMODE", c.SSLMode),
	}
}
