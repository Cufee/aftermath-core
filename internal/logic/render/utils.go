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

func IntToRoman(i int) string {
	roman := ""
	for i > 0 {
		switch {
		case i >= 1000:
			roman += "M"
			i -= 1000
		case i >= 900:
			roman += "CM"
			i -= 900
		case i >= 500:
			roman += "D"
			i -= 500
		case i >= 400:
			roman += "CD"
			i -= 400
		case i >= 100:
			roman += "C"
			i -= 100
		case i >= 90:
			roman += "XC"
			i -= 90
		case i >= 50:
			roman += "L"
			i -= 50
		case i >= 40:
			roman += "XL"
			i -= 40
		case i >= 10:
			roman += "X"
			i -= 10
		case i >= 9:
			roman += "IX"
			i -= 9
		case i >= 5:
			roman += "V"
			i -= 5
		case i >= 4:
			roman += "IV"
			i -= 4
		case i >= 1:
			roman += "I"
			i -= 1
		}
	}
	return roman
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
