package utils

import (
	"regexp"
	"strings"
)

type Slugify struct {
	original string
	slug     string
	pattern  *regexp.Regexp
}

func NewSlugify(s string, pattern string) *Slugify {
	re := regexp.MustCompile(pattern)

	return &Slugify{
		original: s,
		slug:     s,
		pattern:  re,
	}
}

func ToSlug(s string) *Slugify {
	return NewSlugify(s, "[^a-z0-9]+")
}

func (s *Slugify) Slugify() *Slugify {
	slug := strings.ToLower(s.original)
	slug = s.pattern.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	s.slug = slug
	return s
}

func (s *Slugify) AddPrefix(prefix string) *Slugify {
	if s.slug == "" {
		s.Slugify()
	}

	s.slug = strings.TrimRight(prefix, "-") + "-" + s.slug
	return s
}

func (s *Slugify) String() string {
	return s.slug
}
