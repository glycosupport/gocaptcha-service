package store

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"image/color"
	"os"

	"github.com/glycosupport/gocaptcha-service/pkg/captcha"
	"github.com/glycosupport/gocaptcha-service/pkg/utils"
)

type CaptchaData struct {
	Id   string
	Code string
	URL  string
	Hash string
}

type CaptchaStore struct {
	captchas map[string]*CaptchaData
}

type Bg struct {
	R int `json:r`
	G int `json:g`
	B int `json:b`
	A int `json:a`
}

type CaptchaRequest struct {
	Mode   string   `json;type`
	Length int      `json:length`
	Noise  int      `json:noise`
	Lines  int      `json:lines`
	Width  int      `json:width`
	Height int      `json:heigth`
	Fonts  []string `json:fonts`
	Bg     Bg       `json:bg`
	Source string   `json:source`
	Skew   float64  `json:skew`
	Dots   int      `json:dots`
}

func New() *CaptchaStore {
	cs := &CaptchaStore{}
	cs.captchas = make(map[string]*CaptchaData)
	return cs
}

func (cs *CaptchaStore) VerifyCaptcha(hash string, code string) bool {

	if val, ok := cs.captchas[hash]; ok {

		if val.Code == code {
			return true
		} else {
			return false
		}
	}

	return false
}

func (cs *CaptchaStore) GenerateCaptcha(addr string) (*CaptchaData, error) {

	driver := captcha.DefaultDriverDigit

	c := captcha.NewCaptcha(driver, captcha.DefaultMemStore)
	id, b64, answer, res, err := c.Generate()

	if err != nil {
		return nil, err
	}

	hash := utils.GetMD5Hash(b64)

	f, err := os.Create("./captchas/" + hash + ".png")

	if err != nil {
		return nil, err
	}

	defer f.Close()

	_, err = res.WriteTo(f)

	URL := "http://" + addr + "/" + hash + ".png"

	returnValue := &CaptchaData{Id: id, Code: answer, URL: URL, Hash: hash}
	cs.captchas[hash] = returnValue

	return returnValue, err
}

func (cs *CaptchaStore) GenerateCustomCaptcha(addr string, data *CaptchaRequest) (*CaptchaData, error) {

	var driver captcha.Driver

	driver = captcha.NewDriverString(
		data.Height, data.Width, data.Noise,
		data.Lines, data.Length,
		data.Source, &color.RGBA{uint8(data.Bg.R), uint8(data.Bg.G),
			uint8(data.Bg.B), uint8(data.Bg.A)}, nil, data.Fonts)

	switch data.Mode {
	case "string":
		driver = captcha.NewDriverString(data.Height, data.Width, data.Noise, data.Lines, data.Length,
			data.Source, &color.RGBA{uint8(data.Bg.R), uint8(data.Bg.G), uint8(data.Bg.B), uint8(data.Bg.A)}, nil, data.Fonts)
	case "digits":
		driver = captcha.NewDriverDigit(data.Height, data.Width, data.Length, float64(data.Skew), data.Dots)
	case "math":
		driver = captcha.NewDriverMath(data.Height, data.Width, data.Noise, data.Lines,
			&color.RGBA{uint8(data.Bg.R), uint8(data.Bg.G), uint8(data.Bg.B), uint8(data.Bg.A)}, nil, data.Fonts)
	default:
		return nil, errors.New("incorrect mode generate")
	}

	c := captcha.NewCaptcha(driver, captcha.DefaultMemStore)
	id, b64, answer, res, err := c.Generate()

	if err != nil {
		return nil, err
	}

	hash := utils.GetMD5Hash(b64)

	f, err := os.Create("./captchas/" + hash + ".png")

	if err != nil {
		return nil, err
	}

	defer f.Close()

	_, err = res.WriteTo(f)

	URL := "http://" + addr + "/" + hash + ".png"
	returnValue := &CaptchaData{Id: id, Code: answer, URL: URL, Hash: hash}
	cs.captchas[hash] = returnValue

	return returnValue, err
}

func (cs *CaptchaStore) GetCaptcha(name string) (string, error) {

	file, err := os.Open("./captchas/" + name)

	if err != nil {
		return "", err
	}
	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {
		return "", statsErr
	}

	var size int64 = stats.Size()
	bytes := make([]byte, size)

	bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)

	return fmt.Sprintf("data:%s;base64,%s", captcha.MimeTypeImage,
		base64.StdEncoding.EncodeToString(bytes)), err
}
