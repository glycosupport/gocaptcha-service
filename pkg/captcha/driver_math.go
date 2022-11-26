package captcha

import (
	"fmt"
	"image/color"
	"math/rand"
	"strings"

	"github.com/golang/freetype/truetype"
)

type DriverMath struct {
	Height int

	Width int

	NoiseCount int

	ShowLineOptions int

	BgColor *color.RGBA

	fontsStorage FontsStorage

	Fonts      []string
	fontsArray []*truetype.Font
}

func NewDriverMath(height int, width int, noiseCount int, showLineOptions int, bgColor *color.RGBA, fontsStorage FontsStorage, fonts []string) *DriverMath {
	if fontsStorage == nil {
		fontsStorage = DefaultEmbeddedFonts
	}

	tfs := []*truetype.Font{}
	for _, fff := range fonts {
		tf := fontsStorage.LoadFontByName("fonts/" + fff)
		tfs = append(tfs, tf)
	}

	if len(tfs) == 0 {
		tfs = fontsAll
	}

	return &DriverMath{Height: height, Width: width, NoiseCount: noiseCount, ShowLineOptions: showLineOptions, fontsArray: tfs, BgColor: bgColor, Fonts: fonts}
}

func (d *DriverMath) ConvertFonts() *DriverMath {
	if d.fontsStorage == nil {
		d.fontsStorage = DefaultEmbeddedFonts
	}

	tfs := []*truetype.Font{}
	for _, fff := range d.Fonts {
		tf := d.fontsStorage.LoadFontByName("fonts/" + fff)
		tfs = append(tfs, tf)
	}
	if len(tfs) == 0 {
		tfs = fontsAll
	}
	d.fontsArray = tfs

	return d
}

func (d *DriverMath) GenerateIdQuestionAnswer() (id, question, answer string) {
	id = RandomId()
	operators := []string{"+", "-", "x"}
	var mathResult int32
	switch operators[rand.Int31n(3)] {
	case "+":
		a := rand.Int31n(20)
		b := rand.Int31n(20)
		question = fmt.Sprintf("%d+%d=?", a, b)
		mathResult = a + b
	case "x":
		a := rand.Int31n(10)
		b := rand.Int31n(10)
		question = fmt.Sprintf("%dx%d=?", a, b)
		mathResult = a * b
	default:
		a := rand.Int31n(80) + rand.Int31n(20)
		b := rand.Int31n(80)

		question = fmt.Sprintf("%d-%d=?", a, b)
		mathResult = a - b

	}
	answer = fmt.Sprintf("%d", mathResult)
	return
}

func (d *DriverMath) DrawCaptcha(question string) (item Item, err error) {
	var bgc color.RGBA
	if d.BgColor != nil {
		bgc = *d.BgColor
	} else {
		bgc = RandLightColor()
	}
	itemChar := NewItemChar(d.Width, d.Height, bgc)

	if d.ShowLineOptions&OptionShowHollowLine == OptionShowHollowLine {
		itemChar.drawHollowLine()
	}

	if d.NoiseCount > 0 {
		noise := RandText(d.NoiseCount, strings.Repeat(TxtNumbers, d.NoiseCount))
		err = itemChar.drawNoise(noise, fontsAll)
		if err != nil {
			return
		}
	}

	if d.ShowLineOptions&OptionShowSlimeLine == OptionShowSlimeLine {
		itemChar.drawSlimLine(3)
	}

	if d.ShowLineOptions&OptionShowSineLine == OptionShowSineLine {
		itemChar.drawSineLine()
	}

	err = itemChar.drawText(question, d.fontsArray)
	if err != nil {
		return
	}
	return itemChar, nil
}
