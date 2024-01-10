package assets

import (
	"embed"
	"path/filepath"
	"strings"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

//go:embed fonts
var fontsEmbed embed.FS

var Fonts map[string]font.Face

func loadFonts() (map[string]font.Face, error) {
	fontsDir, err := fontsEmbed.ReadDir("fonts")
	if err != nil {
		return nil, err
	}

	fontsMap := make(map[string]font.Face)
	for _, file := range fontsDir {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".ttf") {
			continue
		}
		fontBytes, err := fontsEmbed.ReadFile(filepath.Join("fonts", file.Name()))
		if err != nil {
			return nil, err
		}
		fontFace, err := truetype.Parse(fontBytes)
		if err != nil {
			return nil, err
		}
		face := truetype.NewFace(fontFace, &truetype.Options{
			Size: 16,
		})

		fontsMap[strings.ReplaceAll(file.Name(), ".ttf", "")] = face
		println("loaded font: " + file.Name())
	}

	return fontsMap, nil
}

func init() {
	fonts, err := loadFonts()
	if err != nil {
		panic(err)
	}
	Fonts = fonts
}
