package domain

import (
	"math/rand"

	"github.com/google/uuid"
)

type Avatar struct {
	ID            uuid.UUID
	URL           string
	FallbackChar  string
	FallbackColor string
}

func NewAvatar(name string) *Avatar {
	id := uuid.New()

	url := ""
	fallbackChar, fallbackColor := getFallback(name)

	return &Avatar{
		ID:            id,
		URL:           url,
		FallbackChar:  fallbackChar,
		FallbackColor: fallbackColor,
	}
}

func (a *Avatar) SetURL(url string) {
	a.URL = url
}

func (a *Avatar) SetFallback(name string) {
	a.FallbackChar, a.FallbackColor = getFallback(name)
}

func getFallback(name string) (string, string) {
	return name[0:1], generateRandomColor()
}

func generateRandomColor() string {
	str := "0123456789ABCDEF"
	b := make([]byte, 7)
	b[0] = '#'
	for i := range b[1:] {
		b[i+1] = str[rand.Intn(len(str))]
	}
	return string(b)
}
