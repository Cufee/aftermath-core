package render

import (
	"errors"
	"image"
	"image/color"
	"strings"

	"github.com/EdlinOrg/prominentcolor"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"golang.org/x/image/font"
)

func getDebugColor() color.Color {
	return color.RGBA{255, 192, 203, 255}
}

func GetMatchingColorScheme(img image.Image) ([]color.Color, error) {
	averages, err := prominentcolor.Kmeans(img)
	if err != nil {
		return nil, err
	}
	var colors []color.Color
	for _, average := range averages {
		colors = append(colors, color.RGBA{uint8(average.Color.R), uint8(average.Color.G), uint8(average.Color.B), 120})
	}
	return colors, nil
}

func GetMatchingImageColor(img image.Image) (color.Color, error) {
	matches, err := GetMatchingColorScheme(img)
	if err != nil {
		return nil, err
	}
	if len(matches) == 0 {
		return nil, errors.New("no matching colors found")
	}
	return matches[0], nil
}

func AddBackground(content, background image.Image, style Style) image.Image {
	if background == nil {
		return content
	}

	// Fill background with black and round the corners
	frameCtx := gg.NewContextForImage(content)
	if style.BackgroundColor != nil {
		frameCtx.SetColor(style.BackgroundColor)
		frameCtx.Clear()
	}
	frameCtx.DrawRoundedRectangle(0, 0, float64(frameCtx.Width()), float64(frameCtx.Height()), style.BorderRadius)
	frameCtx.Clip()

	// Resize the background image to fit the cards
	bgImage := imaging.Fill(background, frameCtx.Width(), frameCtx.Height(), imaging.Center, imaging.NearestNeighbor)
	bgImage = imaging.Blur(bgImage, style.Blur)
	frameCtx.DrawImage(bgImage, 0, 0)
	frameCtx.DrawImage(content, 0, 0)

	return frameCtx.Image()
}

type stringSize struct {
	TotalWidth  float64
	TotalHeight float64
	LineOffset  float64
	LineHeight  float64
}

func MeasureString(text string, font font.Face) stringSize {
	if font == nil {
		return stringSize{}
	}

	measureCtx := gg.NewContext(1, 1)
	measureCtx.SetFontFace(font)

	var result stringSize
	// Account for font descender height
	result.LineOffset = float64(font.Metrics().Descent>>6) * 2

	for _, line := range strings.Split(text, "\n") {
		w, h := measureCtx.MeasureString(line)
		h += result.LineOffset
		w += 1

		if w > result.TotalWidth {
			result.TotalWidth = w
		}
		if h > result.LineHeight {
			result.LineHeight = h
		}

		result.TotalHeight += h
	}

	return result
}
