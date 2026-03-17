package service

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type AccountCookieConfig struct {
	AccessTokenCookieConfig  cookieConfig
	RefreshTokenCookieConfig cookieConfig
}

type AccountTokenService struct {
	config AccountCookieConfig
}

func NewAccountTokenService(config AccountCookieConfig) *AccountTokenService {
	return &AccountTokenService{
		config: config,
	}
}

func (s *AccountTokenService) GetRefreshToken(c *fiber.Ctx) (string, error) {
	return c.Cookies(s.config.RefreshTokenCookieConfig.Name), nil
}

func (s *AccountTokenService) GetAccessToken(c *fiber.Ctx) string {
	return c.Cookies(s.config.AccessTokenCookieConfig.Name)
}

func (s *AccountTokenService) Set(c *fiber.Ctx, accessToken, refreshToken string) error {
	c.Cookie(
		&fiber.Cookie{
			Name:     s.config.AccessTokenCookieConfig.Name,
			Value:    accessToken,
			Expires:  s.config.AccessTokenCookieConfig.Expires,
			HTTPOnly: s.config.AccessTokenCookieConfig.HTTPOnly,
			Secure:   s.config.AccessTokenCookieConfig.Secure,
			SameSite: s.config.AccessTokenCookieConfig.SameSite,
		},
	)

	c.Cookie(
		&fiber.Cookie{
			Name:     s.config.RefreshTokenCookieConfig.Name,
			Value:    refreshToken,
			Expires:  s.config.RefreshTokenCookieConfig.Expires,
			HTTPOnly: s.config.RefreshTokenCookieConfig.HTTPOnly,
			Secure:   s.config.RefreshTokenCookieConfig.Secure,
			SameSite: s.config.RefreshTokenCookieConfig.SameSite,
		},
	)

	return nil
}

func (s *AccountTokenService) Clear(c *fiber.Ctx) error {
	c.Cookie(
		&fiber.Cookie{
			Name:     s.config.AccessTokenCookieConfig.Name,
			Value:    "",
			Expires:  time.Now().Add(-24 * time.Hour),
			HTTPOnly: s.config.AccessTokenCookieConfig.HTTPOnly,
			Secure:   s.config.AccessTokenCookieConfig.Secure,
			SameSite: s.config.AccessTokenCookieConfig.SameSite,
		},
	)

	c.Cookie(
		&fiber.Cookie{
			Name:     s.config.RefreshTokenCookieConfig.Name,
			Value:    "",
			Expires:  time.Now().Add(-7 * 24 * time.Hour),
			HTTPOnly: s.config.RefreshTokenCookieConfig.HTTPOnly,
			Secure:   s.config.RefreshTokenCookieConfig.Secure,
			SameSite: s.config.RefreshTokenCookieConfig.SameSite,
		},
	)

	return nil
}
