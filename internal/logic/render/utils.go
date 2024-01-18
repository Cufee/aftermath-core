package render

import (
	"errors"
	"image"
	"image/color"

	"github.com/EdlinOrg/prominentcolor"
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
