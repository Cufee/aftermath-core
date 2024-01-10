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
var fontsMap map[string]*truetype.Font = make(map[string]*truetype.Font)

func init() {
	fonts, err := loadFonts()
	if err != nil {
		panic(err)
	}
	fontsMap = fonts
}
func loadFonts() (map[string]*truetype.Font, error) {
	fontsDir, err := fontsEmbed.ReadDir("fonts")
	if err != nil {
		return nil, err
	}

	fontsMap := make(map[string]*truetype.Font)
	for _, file := range fontsDir {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".ttf") {
			continue
		}
		fontBytes, err := fontsEmbed.ReadFile(filepath.Join("fonts", file.Name()))
		if err != nil {
			return nil, err
		}
		font, err := truetype.Parse(fontBytes)
		if err != nil {
			return nil, err
		}
		fontsMap[strings.ReplaceAll(file.Name(), ".ttf", "")] = font
		println("loaded font: " + file.Name())
	}

	return fontsMap, nil
}

func GetFontFaces(name string, sizes ...float64) (map[float64]font.Face, bool) {
	loadedFont, ok := fontsMap[name]
	if !ok {
		return nil, false
	}
	faces := make(map[float64]font.Face)
	for _, size := range sizes {
		faces[size] = truetype.NewFace(loadedFont, &truetype.Options{
			Size: size,
		})
	}
	return faces, true
}
