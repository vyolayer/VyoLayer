package service

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type IAMCookie string

const (
	IAMCookieAccessToken  IAMCookie = "__vyo_iam_auth"
	IAMCookieRefreshToken IAMCookie = "__vyo_iam_session"
)

type IAMCookieConfig struct {
	Atcc cookieConfig
	Rtcc cookieConfig
}

type IAMCookieService struct {
	config IAMCookieConfig
}

func NewIAMCookieService(config IAMCookieConfig) *IAMCookieService {
	return &IAMCookieService{
		config: config,
	}
}

// Get cookies
func (s *IAMCookieService) GetAccessToken(c *fiber.Ctx) string {
	return c.Cookies(s.config.Atcc.Name)
}

func (s *IAMCookieService) GetSessionToken(c *fiber.Ctx) string {
	return c.Cookies(s.config.Rtcc.Name)
}

// Add cookies
func (s *IAMCookieService) Set(c *fiber.Ctx, accessToken, refreshToken string, atccExpiryAt, rtccExpiryAt time.Time) error {

	c.Cookie(
		&fiber.Cookie{
			Name:     s.config.Atcc.Name,
			Value:    accessToken,
			MaxAge:   int(time.Until(atccExpiryAt).Seconds()),
			HTTPOnly: s.config.Atcc.HTTPOnly,
			Secure:   s.config.Atcc.Secure,
			SameSite: s.config.Atcc.SameSite,
		},
	)

	c.Cookie(
		&fiber.Cookie{
			Name:     s.config.Rtcc.Name,
			Value:    refreshToken,
			MaxAge:   int(time.Until(rtccExpiryAt).Seconds()),
			HTTPOnly: s.config.Rtcc.HTTPOnly,
			Secure:   s.config.Rtcc.Secure,
			SameSite: s.config.Rtcc.SameSite,
		},
	)

	return nil
}

func (s *IAMCookieService) Clear(c *fiber.Ctx) error {
	c.Cookie(
		&fiber.Cookie{
			Name:     s.config.Atcc.Name,
			Value:    "",
			Expires:  time.Now().Add(-24 * time.Hour),
			HTTPOnly: s.config.Atcc.HTTPOnly,
			Secure:   s.config.Atcc.Secure,
			SameSite: s.config.Atcc.SameSite,
		},
	)

	c.Cookie(
		&fiber.Cookie{
			Name:     s.config.Rtcc.Name,
			Value:    "",
			Expires:  time.Now().Add(-7 * 24 * time.Hour),
			HTTPOnly: s.config.Rtcc.HTTPOnly,
			Secure:   s.config.Rtcc.Secure,
			SameSite: s.config.Rtcc.SameSite,
		},
	)

	return nil
}
