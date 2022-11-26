package captcha

import (
	"embed"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

type EmbeddedFontsStorage struct {
	fs embed.FS
}

func (s *EmbeddedFontsStorage) LoadFontByName(name string) *truetype.Font {
	fontBytes, err := s.fs.ReadFile(name)
	if err != nil {
		panic(err)
	}

	trueTypeFont, err := freetype.ParseFont(fontBytes)
	if err != nil {
		panic(err)
	}

	return trueTypeFont
}

func (s *EmbeddedFontsStorage) LoadFontsByNames(assetFontNames []string) []*truetype.Font {
	fonts := make([]*truetype.Font, 0)
	for _, assetName := range assetFontNames {
		f := s.LoadFontByName(assetName)
		fonts = append(fonts, f)
	}
	return fonts
}

func NewEmbeddedFontsStorage(fs embed.FS) *EmbeddedFontsStorage {
	return &EmbeddedFontsStorage{
		fs: fs,
	}
}
