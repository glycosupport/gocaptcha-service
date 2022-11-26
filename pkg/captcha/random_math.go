package captcha

import (
	"encoding/binary"
	"image/color"
	"math"
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandText(size int, sourceChars string) string {
	if sourceChars == "" || size == 0 {
		return ""
	}

	if size >= len(sourceChars) {
		sourceChars = strings.Repeat(sourceChars, size)
	}

	sourceRunes := []rune(sourceChars)
	sourceLength := len(sourceRunes)

	text := make([]rune, size)
	for i := range text {
		text[i] = sourceRunes[rand.Intn(sourceLength)]
	}
	return string(text)
}

func random(min int64, max int64) float64 {
	return float64(min) + rand.Float64()*float64(max-min)
}

func RandDeepColor() color.RGBA {

	randColor := RandColor()

	increase := float64(30 + rand.Intn(255))

	red := math.Abs(math.Min(float64(randColor.R)-increase, 255))

	green := math.Abs(math.Min(float64(randColor.G)-increase, 255))
	blue := math.Abs(math.Min(float64(randColor.B)-increase, 255))

	return color.RGBA{R: uint8(red), G: uint8(green), B: uint8(blue), A: uint8(255)}
}

func RandLightColor() color.RGBA {
	red := rand.Intn(55) + 200
	green := rand.Intn(55) + 200
	blue := rand.Intn(55) + 200
	return color.RGBA{R: uint8(red), G: uint8(green), B: uint8(blue), A: uint8(255)}
}

func RandColor() color.RGBA {
	red := rand.Intn(255)
	green := rand.Intn(255)
	var blue int
	if (red + green) > 400 {
		blue = 0
	} else {
		blue = 400 - green - red
	}
	if blue > 255 {
		blue = 255
	}
	return color.RGBA{R: uint8(red), G: uint8(green), B: uint8(blue), A: uint8(255)}
}

func randIntRange(from, to int) int {
	if to-from <= 0 {
		return from
	}
	return rand.Intn(to-from) + from
}
func randFloat64Range(from, to float64) float64 {
	return rand.Float64()*(to-from) + from
}
func randBytes(n int) []byte {
	numBlocks := (n + 8 - 1) / 8
	b := make([]byte, numBlocks*8)
	for i := 0; i < len(b); i += 8 {

		binary.LittleEndian.PutUint64(b[i:], rand.Uint64())
	}
	return b[:n]
}

func RandomId() string {
	b := randomBytesMod(idLen, byte(len(idChars)))
	for i, c := range b {
		b[i] = idChars[c]
	}
	return string(b)
}
