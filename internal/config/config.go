package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Auth     AuthConfig     `yaml:"auth"`
}

type AppConfig struct {
	Name string `yaml:"name"`
	Env  string `yaml:"environment"` // development, production, test
}

type ServerConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

type DatabaseConfig struct {
	Driver          string        `yaml:"driver"`
	DSN             string        `yaml:"dsn"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

type AuthConfig struct {
	JWTSecret          string        `yaml:"jwt_secret"`
	AccessTokenTTL     time.Duration `yaml:"access_token_ttl"`
	RefreshTokenSecret string        `yaml:"refresh_token_secret"`
	RefreshTokenTTL    time.Duration `yaml:"refresh_token_ttl"`
}

func Load(path string) (*Config, error) {

	// Load .env if present (non-fatal)
	_ = godotenv.Load()

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	applyEnvOverrides(&cfg)

	if err := validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func applyEnvOverrides(cfg *Config) {

	if v := os.Getenv("APP_ENV"); v != "" {
		cfg.App.Env = v
	}

	if v := os.Getenv("SERVER_PORT"); v != "" {
		fmt.Sscan(v, &cfg.Server.Port)
	}

	if v := os.Getenv("DATABASE_DSN"); v != "" {
		cfg.Database.DSN = v
	}

	if v := os.Getenv("JWT_SECRET"); v != "" {
		cfg.Auth.JWTSecret = v
	}
}

func validate(cfg *Config) error {

	if cfg.App.Name == "" {
		return errors.New("app.name is required")
	}

	if cfg.Server.Port == 0 {
		return errors.New("server.port is required")
	}

	if cfg.Database.DSN == "" {
		return errors.New("database.dsn is required")
	}

	if cfg.Auth.JWTSecret == "" {
		return errors.New("auth.jwt_secret is required")
	}

	return nil
}
