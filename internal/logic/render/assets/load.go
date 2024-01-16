package assets

import (
	"bytes"
	"embed"
	"image"
	"path/filepath"
	"strings"

	"github.com/golang/freetype/truetype"
	"github.com/rs/zerolog/log"
	"golang.org/x/image/font"
)

//go:embed fonts
var fontsEmbed embed.FS

//go:embed images
var imagesEmbed embed.FS

var fontsMap map[string]*truetype.Font = make(map[string]*truetype.Font)
var imagesMap map[string]image.Image = make(map[string]image.Image)

func init() {
	fonts, err := loadFonts()
	if err != nil {
		panic(err)
	}
	fontsMap = fonts

	images, err := loadImages()
	if err != nil {
		panic(err)
	}
	imagesMap = images
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
		log.Debug().Msg("loaded font: " + file.Name())
	}

	return fontsMap, nil
}

func loadImages() (map[string]image.Image, error) {
	images, err := getAllFiles(imagesEmbed, ".")
	if err != nil {
		return nil, err
	}

	for path, file := range images {
		if (!strings.HasSuffix(path, ".png")) && !strings.HasSuffix(path, ".jpg") {
			continue
		}

		image, _, err := image.Decode(bytes.NewBuffer(file))
		if err != nil {
			return nil, err
		}

		imagesMap[strings.Split(path, ".")[0]] = image
		log.Debug().Msg("loaded image: " + path)
	}

	return imagesMap, nil
}

func getAllFiles(dir embed.FS, path string) (map[string][]byte, error) {
	entries, err := dir.ReadDir(path)
	if err != nil {
		return nil, err
	}

	files := make(map[string][]byte)
	for _, entry := range entries {
		if entry.IsDir() {
			subFiles, err := getAllFiles(dir, filepath.Join(path, entry.Name()))
			if err != nil {
				return nil, err
			}
			for k, v := range subFiles {
				files[k] = v
			}
			continue
		}

		file, err := dir.ReadFile(filepath.Join(path, entry.Name()))
		if err != nil {
			return nil, err
		}
		files[filepath.Join(path, entry.Name())] = file
	}
	return files, nil
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

func GetImage(name string) (image.Image, bool) {
	img, ok := imagesMap[name]
	return img, ok
}
