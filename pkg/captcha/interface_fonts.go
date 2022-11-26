package captcha

import "github.com/golang/freetype/truetype"

type FontsStorage interface {
	LoadFontByName(name string) *truetype.Font

	LoadFontsByNames(assetFontNames []string) []*truetype.Font
}
