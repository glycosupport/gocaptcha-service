package captcha

type Driver interface {
	DrawCaptcha(content string) (item Item, err error)
	GenerateIdQuestionAnswer() (id, q, a string)
}
