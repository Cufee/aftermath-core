package render

import (
	"errors"
	"image"
	"image/color"

	"github.com/EdlinOrg/prominentcolor"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
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
