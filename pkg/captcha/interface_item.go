package captcha

import "io"

type Item interface {
	WriteTo(w io.Writer) (n int64, err error)
	EncodeB64string() string
}
