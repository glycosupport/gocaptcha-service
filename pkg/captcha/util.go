package captcha

import (
	"crypto/rand"
	"io"
)

func parseDigitsToString(bytes []byte) string {
	stringB := make([]byte, len(bytes))
	for idx, by := range bytes {
		stringB[idx] = by + '0'
	}
	return string(stringB)
}

func stringToFakeByte(content string) []byte {
	digits := make([]byte, len(content))
	for idx, cc := range content {
		digits[idx] = byte(cc - '0')
	}
	return digits
}

func randomDigits(length int) []byte {
	return randomBytesMod(length, 10)
}

func randomBytes(length int) (b []byte) {
	b = make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic("captcha: error reading random source: " + err.Error())
	}
	return
}

func randomBytesMod(length int, mod byte) (b []byte) {
	if length == 0 {
		return nil
	}
	if mod == 0 {
		panic("captcha: bad mod argument for randomBytesMod")
	}
	maxrb := 255 - byte(256%int(mod))
	b = make([]byte, length)
	i := 0
	for {
		r := randomBytes(length + (length / 4))
		for _, c := range r {
			if c > maxrb {
				continue
			}
			b[i] = c % mod
			i++
			if i == length {
				return
			}
		}
	}
}
