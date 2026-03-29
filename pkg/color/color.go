package color

import "math/rand"

func GenerateColor() string {
	str := "0123456789ABCDEF"
	b := make([]byte, 7)
	b[0] = '#'
	for i := range b[1:] {
		b[i+1] = str[rand.Intn(len(str))]
	}
	return string(b)
}
