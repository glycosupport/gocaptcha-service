package captcha

// Captcha основная стркутура пакета
type Captcha struct {
	Driver Driver
	Store  Store
}

// NewCaptcha генерирует новую каптчу
func NewCaptcha(driver Driver, store Store) *Captcha {
	return &Captcha{Driver: driver, Store: store}
}

// Generate генерирует рандомный id, base64 изоображение и item для просмотра сгенерированного текста
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
