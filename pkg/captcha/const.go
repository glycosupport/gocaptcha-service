package captcha

const idLen = 20

var idChars = []byte(TxtNumbers + TxtAlphabet)

const (
	imageStringDpi     = 72.0
	TxtNumbers         = "012346789"
	TxtAlphabet        = "ABCDEFGHJKMNOQRSTUVXYZabcdefghjkmnoqrstuvxyz"
	TxtSimpleCharaters = "13467ertyiadfhjkxcvbnERTYADFGHJKXCVBN"
	MimeTypeImage      = "image/png"
)

const (
	OptionShowHollowLine = 2
	OptionShowSlimeLine  = 4
	OptionShowSineLine   = 8
)

const (
	DefaultHeight     = 60
	DefaultWidth      = 240
	DefaultNoiseCount = 0
	DefaultShowLine   = 3
	DefaultLength     = 6
	DefaultSource     = TxtAlphabet
)
