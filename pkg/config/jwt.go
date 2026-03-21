package config

import "time"

type JWTConfig struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

var (
	DefaultJWTConfig = JWTConfig{
		AccessTokenSecret:  "access_token_secret",
		RefreshTokenSecret: "refresh_token_secret",
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenExpiry: 168 * time.Hour,
	}
)

func NewJWTConfig(c JWTConfig) *JWTConfig {
	return &JWTConfig{
		AccessTokenSecret:  GetEnv("ACCESS_TOKEN_SECRET", c.AccessTokenSecret),
		AccessTokenExpiry:  GetEnvDuration("ACCESS_TOKEN_EXPIRY", c.AccessTokenExpiry.String()),
		RefreshTokenSecret: GetEnv("REFRESH_TOKEN_SECRET", c.RefreshTokenSecret),
		RefreshTokenExpiry: GetEnvDuration("REFRESH_TOKEN_EXPIRY", c.RefreshTokenExpiry.String()),
	}
}
