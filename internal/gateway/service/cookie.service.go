package service

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type CookieService interface {
	Set(c *fiber.Ctx, token string) error
	Get(c *fiber.Ctx) string
	Clear(c *fiber.Ctx) error
}

type cookieConfig = fiber.Cookie

type cookieService struct {
	config cookieConfig
}

func NewCookieService(config cookieConfig) CookieService {
	return &cookieService{
		config: config,
	}
}

func (s *cookieService) Get(c *fiber.Ctx) string {
	return c.Cookies(s.config.Name)
}

func (s *cookieService) Set(c *fiber.Ctx, token string) error {
	c.Cookie(
		&fiber.Cookie{
			Name:     s.config.Name,
			Value:    token,
			Expires:  s.config.Expires,
			HTTPOnly: s.config.HTTPOnly,
			Secure:   s.config.Secure,
			SameSite: s.config.SameSite,
		},
	)

	return nil
}

func (s *cookieService) Clear(c *fiber.Ctx) error {
	c.Cookie(
		&fiber.Cookie{
			Name:     s.config.Name,
			Value:    "",
			Expires:  time.Now().Add(-24 * time.Hour),
			HTTPOnly: s.config.HTTPOnly,
			Secure:   s.config.Secure,
			SameSite: s.config.SameSite,
		},
	)

	return nil
}
