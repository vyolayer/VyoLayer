package domain

import (
	"github.com/vyolayer/vyolayer/pkg/color"
)

type IAMUserAvatar struct {
	ID            int64
	URL           string
	FallbackChar  string
	FallbackColor string
}

func NewIAMUserAvatar(URL, fallbackChar string) *IAMUserAvatar {
	return &IAMUserAvatar{
		URL:           URL,
		FallbackChar:  fallbackChar,
		FallbackColor: color.GenerateColor(),
	}
}
