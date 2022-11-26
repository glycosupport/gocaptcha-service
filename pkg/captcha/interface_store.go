package captcha

type Store interface {
	Set(id string, value string) error

	Get(id string, clear bool) string

	Verify(id, answer string, clear bool) bool
}
