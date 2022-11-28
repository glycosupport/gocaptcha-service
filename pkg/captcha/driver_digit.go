package captcha

import "math/rand"

type DriverDigit struct {
	Height   int
	Width    int
	Length   int
	MaxSkew  float64
	DotCount int
}

func NewDriverDigit(height int, width int, length int, maxSkew float64, dotCount int) *DriverDigit {
	return &DriverDigit{Height: height, Width: width, Length: length, MaxSkew: maxSkew, DotCount: dotCount}
}

var DefaultDriverDigit = NewDriverDigit(80, 240, 6, 0.7, 80)

func (d *DriverDigit) GenerateIdQuestionAnswer() (id, q, a string) {
	id = RandomId()
	digits := randomDigits(d.Length)
	a = parseDigitsToString(digits)
	return id, a, a
}

func (d *DriverDigit) DrawCaptcha(content string) (item Item, err error) {
	itemDigit := NewItemDigit(d.Width, d.Height, d.DotCount, d.MaxSkew)
	digits := stringToFakeByte(content)

	itemDigit.calculateSizes(d.Width, d.Height, len(digits))
	maxx := d.Width - (itemDigit.width+itemDigit.dotSize)*len(digits) - itemDigit.dotSize
	maxy := d.Height - itemDigit.height - itemDigit.dotSize*2

	var border int

	if d.Width > d.Height {
		border = d.Height / 5
	} else {
		border = d.Width / 5
	}

	x := rand.Intn(maxx-border*2) + border
	y := rand.Intn(maxy-border*2) + border

	for _, n := range digits {
		itemDigit.drawDigit(digitFontData[n], x, y)
		x += itemDigit.width + itemDigit.dotSize
	}

	itemDigit.strikeThrough()
	itemDigit.distort(rand.Float64()*(10-5)+5, rand.Float64()*(200-100)+100)
	itemDigit.fillWithCircles(d.DotCount, itemDigit.dotSize)

	return itemDigit, nil
}
