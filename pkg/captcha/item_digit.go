package captcha

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"math/rand"
)

const (
	digitFontWidth     = 11
	digitFontHeight    = 18
	digitFontBlackChar = 1
)

type ItemDigit struct {
	width  int
	height int
	*image.Paletted
	dotSize  int
	dotCount int
	maxSkew  float64
}

func NewItemDigit(width int, height int, dotCount int, maxSkew float64) *ItemDigit {
	itemDigit := &ItemDigit{width: width, height: height, dotCount: dotCount, maxSkew: maxSkew}
	itemDigit.Paletted = image.NewPaletted(image.Rect(0, 0, width, height), createRandPaletteColors(dotCount))
	return itemDigit
}

func createRandPaletteColors(dotCount int) color.Palette {
	p := make([]color.Color, dotCount+1)
	p[0] = color.RGBA{0xFF, 0xFF, 0xFF, 0x00}
	prim := color.RGBA{
		uint8(rand.Intn(129)),
		uint8(rand.Intn(129)),
		uint8(rand.Intn(129)),
		0xFF,
	}

	if dotCount == 0 {
		p[0] = prim
		return p
	}

	p[1] = prim
	for i := 2; i <= dotCount; i++ {
		p[i] = randomBrightness(prim, 255)
	}
	return p
}

func (m *ItemDigit) calculateSizes(width, height, ncount int) {
	var border int
	if width > height {
		border = height / 4
	} else {
		border = width / 4
	}

	w := float64(width - border*2)
	h := float64(height - border*2)

	fw := float64(digitFontWidth + 1)
	fh := float64(digitFontHeight)
	nc := float64(ncount)

	nw := w / nc

	nh := nw * fh / fw

	if nh > h {
		nh = h
		nw = fw / fh * nh
	}

	m.dotSize = int(nh / fh)
	if m.dotSize < 1 {
		m.dotSize = 1
	}

	m.width = int(nw) - m.dotSize
	m.height = int(nh)
}

func (m *ItemDigit) drawHorizLine(fromX, toX, y int, colorIdx uint8) {
	for x := fromX; x <= toX; x++ {
		m.SetColorIndex(x, y, colorIdx)
	}
}

func (m *ItemDigit) drawCircle(x, y, radius int, colorIdx uint8) {
	f := 1 - radius
	dfx := 1
	dfy := -2 * radius
	xo := 0
	yo := radius

	m.SetColorIndex(x, y+radius, colorIdx)
	m.SetColorIndex(x, y-radius, colorIdx)
	m.drawHorizLine(x-radius, x+radius, y, colorIdx)

	for xo < yo {
		if f >= 0 {
			yo--
			dfy += 2
			f += dfy
		}
		xo++
		dfx += 2
		f += dfx
		m.drawHorizLine(x-xo, x+xo, y+yo, colorIdx)
		m.drawHorizLine(x-xo, x+xo, y-yo, colorIdx)
		m.drawHorizLine(x-yo, x+yo, y+xo, colorIdx)
		m.drawHorizLine(x-yo, x+yo, y-xo, colorIdx)
	}
}

func (m *ItemDigit) fillWithCircles(n, maxradius int) {
	maxx := m.Bounds().Max.X
	maxy := m.Bounds().Max.Y
	for i := 0; i < n; i++ {
		colorIdx := uint8(randIntRange(1, m.dotCount-1))
		r := randIntRange(1, maxradius)
		m.drawCircle(randIntRange(r, maxx-r), randIntRange(r, maxy-r), r, colorIdx)
	}
}

func (m *ItemDigit) strikeThrough() {
	maxx := m.Bounds().Max.X
	maxy := m.Bounds().Max.Y
	y := randIntRange(maxy/3, maxy-maxy/3)
	amplitude := randFloat64Range(5, 20)
	period := randFloat64Range(80, 180)
	dx := 2.0 * math.Pi / period
	for x := 0; x < maxx; x++ {
		xo := amplitude * math.Cos(float64(y)*dx)
		yo := amplitude * math.Sin(float64(x)*dx)
		for yn := 0; yn < m.dotSize; yn++ {
			r := rand.Intn(m.dotSize)
			m.drawCircle(x+int(xo), y+int(yo)+(yn*m.dotSize), r/2, 1)
		}
	}
}

func (m *ItemDigit) drawDigit(digit []byte, x, y int) {
	skf := randFloat64Range(-m.maxSkew, m.maxSkew)
	xs := float64(x)
	r := m.dotSize / 2
	y += randIntRange(-r, r)
	for yo := 0; yo < digitFontHeight; yo++ {
		for xo := 0; xo < digitFontWidth; xo++ {
			if digit[yo*digitFontWidth+xo] != digitFontBlackChar {
				continue
			}
			m.drawCircle(x+xo*m.dotSize, y+yo*m.dotSize, r, 1)
		}
		xs += skf
		x = int(xs)
	}
}

func (m *ItemDigit) distort(amplude float64, period float64) {
	w := m.Bounds().Max.X
	h := m.Bounds().Max.Y

	oldm := m.Paletted
	newm := image.NewPaletted(image.Rect(0, 0, w, h), oldm.Palette)

	dx := 2.0 * math.Pi / period
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			xo := amplude * math.Sin(float64(y)*dx)
			yo := amplude * math.Cos(float64(x)*dx)
			newm.SetColorIndex(x, y, oldm.ColorIndexAt(x+int(xo), y+int(yo)))
		}
	}
	m.Paletted = newm
}

func randomBrightness(c color.RGBA, max uint8) color.RGBA {
	minc := min3(c.R, c.G, c.B)
	maxc := max3(c.R, c.G, c.B)
	if maxc > max {
		return c
	}
	n := rand.Intn(int(max-maxc)) - int(minc)
	return color.RGBA{
		uint8(int(c.R) + n),
		uint8(int(c.G) + n),
		uint8(int(c.B) + n),
		uint8(c.A),
	}
}

func min3(x, y, z uint8) (m uint8) {
	m = x
	if y < m {
		m = y
	}
	if z < m {
		m = z
	}
	return
}

func max3(x, y, z uint8) (m uint8) {
	m = x
	if y > m {
		m = y
	}
	if z > m {
		m = z
	}
	return
}

func (m *ItemDigit) EncodeBinary() []byte {
	var buf bytes.Buffer
	if err := png.Encode(&buf, m.Paletted); err != nil {
		panic(err.Error())
	}
	return buf.Bytes()
}

func (m *ItemDigit) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(m.EncodeBinary())
	return int64(n), err
}

func (m *ItemDigit) EncodeB64string() string {
	return fmt.Sprintf("data:%s;base64,%s", MimeTypeImage, base64.StdEncoding.EncodeToString(m.EncodeBinary()))
}
