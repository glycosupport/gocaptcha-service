package captcha

type Captcha struct {
	Driver Driver
}

func NewCaptcha(driver Driver) *Captcha {
	return &Captcha{Driver: driver}
}

func (c *Captcha) Generate() (id, b64, answer string, item Item, err error) {
	id, content, answer := c.Driver.GenerateIdQuestionAnswer()
	item, err = c.Driver.DrawCaptcha(content)

	if err != nil {
		return "", "", "", nil, err
	}

	if err != nil {
		return "", "", "", nil, err
	}

	b64 = item.EncodeB64string()
	return
}
