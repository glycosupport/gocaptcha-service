package captcha

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/golang/freetype/truetype"
)

type DriverString struct {
	Height int

	Width int

	NoiseCount int

	ShowLineOptions int

	Length int

	Source string

	BgColor *color.RGBA

	fontsStorage FontsStorage

	Fonts      []string
	fontsArray []*truetype.Font
}

func NewDriverString(height int, width int, noiseCount int, showLineOptions int, length int, source string, bgColor *color.RGBA, fontsStorage FontsStorage, fonts []string) *DriverString {
	if fontsStorage == nil {
		fontsStorage = DefaultEmbeddedFonts
	}

	tfs := []*truetype.Font{}
	for _, fff := range fonts {
		tf := fontsStorage.LoadFontByName(fff)
		tfs = append(tfs, tf)
	}

	if len(tfs) == 0 {
		tfs = fontsAll
	}

	return &DriverString{Height: height, Width: width, NoiseCount: noiseCount, ShowLineOptions: showLineOptions, Length: length, Source: source, BgColor: bgColor, fontsStorage: fontsStorage, fontsArray: tfs, Fonts: fonts}
}

func (d *DriverString) ConvertFonts() *DriverString {
	if d.fontsStorage == nil {
		d.fontsStorage = DefaultEmbeddedFonts
	}

	tfs := []*truetype.Font{}
	for _, fff := range d.Fonts {
		tf := d.fontsStorage.LoadFontByName(fff)
		tfs = append(tfs, tf)
	}
	if len(tfs) == 0 {
		tfs = fontsAll
	}

	d.fontsArray = tfs

	return d
}

func (d *DriverString) GenerateIdQuestionAnswer() (id, content, answer string) {
	id = RandomId()
	content = strings.ToLower(RandText(d.Length, d.Source))
	return id, content, content
}

func (d *DriverString) DrawCaptcha(content string) (item Item, err error) {

	var bgc color.RGBA

	fmt.Println(*d.BgColor)

	if d.BgColor != nil {
		bgc = *d.BgColor
	} else {
		bgc = RandLightColor()
	}
	itemChar := NewItemChar(d.Width, d.Height, bgc)

	if d.ShowLineOptions&OptionShowHollowLine == OptionShowHollowLine {
		itemChar.drawHollowLine()
	}

	if d.ShowLineOptions&OptionShowSlimeLine == OptionShowSlimeLine {
		itemChar.drawSlimLine(3)
	}

	if d.ShowLineOptions&OptionShowSineLine == OptionShowSineLine {
		itemChar.drawSineLine()
	}

	if d.NoiseCount > 0 {
		source := TxtNumbers + TxtAlphabet + ",.[]<>"
		noise := RandText(d.NoiseCount, strings.Repeat(source, d.NoiseCount))
		err = itemChar.drawNoise(noise, d.fontsArray)
		if err != nil {
			return
		}
	}

	err = itemChar.drawText(content, d.fontsArray)
	if err != nil {
		return
	}

	return itemChar, nil
}
